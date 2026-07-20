package networking

import (
	"fmt"
)

type Command interface {
	execute(cm *ConnectionManager, rm *RoomManager, senderID int) error
}

type moveCommand struct {
	move string // LEFT, RIGHT, DOWN, DROP
}

type createRoomCommand struct {
	roomID int
}

type joinRoomCommand struct {
	roomID int
}

func (cmd *moveCommand) execute(cm *ConnectionManager, rm *RoomManager, senderID int) error {
	roomID, ok := rm.GetRoom(senderID)
	if !ok {
		return fmt.Errorf("not in a room")
	}
	message := fmt.Sprintf("%d %s", senderID, cmd.move)
	return rm.RoomBroadcastExcept(roomID, senderID, message, cm)
}

func (cmd *createRoomCommand) execute(cm *ConnectionManager, rm *RoomManager, senderID int) error {
	if rm.RoomExists(cmd.roomID) {
		return fmt.Errorf("room %d already exists", cmd.roomID)
	}

	fmt.Println("info : ", cmd.roomID, senderID)
	rm.Join(cmd.roomID, senderID)
	message := fmt.Sprintf("JOINED %d %d", cmd.roomID, senderID)
	return rm.RoomBroadcast(cmd.roomID, message, cm)
}

func (cmd *joinRoomCommand) execute(cm *ConnectionManager, rm *RoomManager, senderID int) error {
	if !rm.RoomExists(cmd.roomID) {
		return fmt.Errorf("room %d does not exist", cmd.roomID)
	}
	if rm.ClientInRoom(cmd.roomID, senderID) {
		return fmt.Errorf("already in room %d", cmd.roomID)
	}
	rm.Join(cmd.roomID, senderID)
	message := fmt.Sprintf("JOINED %d %d", cmd.roomID, senderID)
	return rm.RoomBroadcast(cmd.roomID, message, cm)
}
