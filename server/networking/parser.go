package networking

import (
	"strconv"
	"strings"
)

func parseCommand(line string) Command {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}

	switch fields[0] {
	case "LEFT", "RIGHT", "DOWN", "DROP":
		return &moveCommand{move: fields[0]}

	case "CREATE_ROOM":
		if len(fields) < 2 {
			return nil
		}
		roomID, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil
		}
		return &createRoomCommand{roomID: roomID}

	case "JOIN_ROOM":
		if len(fields) < 2 {
			return nil
		}
		roomID, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil
		}
		return &joinRoomCommand{roomID: roomID}
	}

	return nil
}
