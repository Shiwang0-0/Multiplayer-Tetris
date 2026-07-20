package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/Shiwang0-0/multiplayertetris/protocol"
)

func ParseResponse(line string) tea.Msg {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	space := strings.IndexByte(line, ' ')
	var head, rest string
	if space == -1 {
		head = line
	} else {
		head = line[:space]
		rest = line[space+1:]
	}
	switch head {
	case "CLIENT_ID":
		return parseClientID(rest)
	case "JOINED":
		return parseJoined(rest)
	case "ERROR":
		return parseError(rest)
	default:
		// Not a known keyword — this is "<senderID> <MOVE>".
		return parseOpponentMove(head, rest)
	}
}

func parseClientID(data string) tea.Msg {
	id, err := strconv.Atoi(data)
	if err != nil {
		return nil
	}
	return protocol.ClientIDMsg{ID: id}
}

func parseJoined(data string) tea.Msg { // JOINED Room: 3 - Client: 2
	var roomID int
	var clientID int

	_, err := fmt.Sscanf(data, "%d %d", &roomID, &clientID)

	if err != nil {
		log.Println("Invalid JOINED response")
		return nil
	}

	return protocol.JoinedMsg{
		RoomID:   roomID,
		ClientID: clientID,
	}
}

func parseOpponentMove(senderIDStr, move string) tea.Msg {
	senderID, err := strconv.Atoi(senderIDStr)
	if err != nil {
		return nil
	}
	move = strings.TrimSpace(move)
	return protocol.OpponentMoveMsg{
		SenderID: senderID,
		Move:     move,
	}
}

func parseError(data string) tea.Msg {
	return protocol.ServerErrorMsg{Msg: data}
}
