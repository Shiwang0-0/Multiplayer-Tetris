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

type voteCommand struct {
	piece string
}

type lockedCommand struct{}

type startMatchCommand struct{}

type playerOutCommand struct{}

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
	rm.Join(cmd.roomID, senderID)

	rm.mu.Lock()
	rm.roomCreator[cmd.roomID] = senderID
	rm.mu.Unlock()

	message := fmt.Sprintf("JOINED %d %d", cmd.roomID, senderID)
	if err := rm.RoomBroadcast(cmd.roomID, message, cm); err != nil {
		return err
	}
	return nil
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

// after room waiting finidh and game starts
func (cmd *voteCommand) execute(cm *ConnectionManager, rm *RoomManager, senderID int) error {
	roomID, ok := rm.GetRoom(senderID)
	if !ok {
		return fmt.Errorf("not in a room")
	}
	match, ok := rm.GetMatch(roomID)
	if !ok {
		return fmt.Errorf("no active match")
	}
	return match.SubmitVote(senderID, cmd.piece)
}

// after one piece is locked,
func (cmd *lockedCommand) execute(cm *ConnectionManager, rm *RoomManager, senderID int) error {
	roomID, ok := rm.GetRoom(senderID)
	if !ok {
		return fmt.Errorf("not in a room")
	}
	match, ok := rm.GetMatch(roomID)
	if !ok {
		return fmt.Errorf("no active match")
	}
	return match.OnLocked(senderID, cm, rm)
}

func (cmd *startMatchCommand) execute(cm *ConnectionManager, rm *RoomManager, senderID int) error {
	roomID, ok := rm.GetRoom(senderID)
	if !ok {
		return fmt.Errorf("not in a room")
	}
	if !rm.IsCreator(roomID, senderID) {
		return fmt.Errorf("only the room creator can start the match")
	}
	if rm.RoomSize(roomID) < 2 {
		return fmt.Errorf("waiting for more players")
	}
	rm.StartMatchIfReady(roomID, cm)
	return nil
}

func (cmd *playerOutCommand) execute(cm *ConnectionManager, rm *RoomManager, senderID int) error {
	roomID, ok := rm.GetRoom(senderID)
	if !ok {
		return fmt.Errorf("not in a room")
	}
	match, ok := rm.GetMatch(roomID)
	if !ok {
		return fmt.Errorf("no active match")
	}
	return match.OnPlayerOut(senderID, cm, rm)
}
