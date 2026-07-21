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
	case "VOTING_START":
		return parseVotingStart(rest)
	case "TURN_START":
		return parseTurnStart(rest)
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

func parseVotingStart(data string) tea.Msg {

	var activePlayerID int
	var deadline int

	_, err := fmt.Sscanf(data, "%d %d", &activePlayerID, &deadline)
	if err != nil {
		log.Println("Invalid VOITING_START response")
		return nil
	}
	return protocol.VotingStartMsg{
		ActivePlayerID:  activePlayerID,
		DeadlineSeconds: deadline,
	}
}

func parseTurnStart(data string) tea.Msg {
	var activePlayerID int // who the active playing player is (who is getting the votes)
	var winnerPiece string

	_, err := fmt.Sscanf(data, "%d %s", &activePlayerID, &winnerPiece)
	if err != nil {
		log.Println("Invalid VOITING_START response")
		return nil
	}
	return protocol.TurnStartMsg{
		ActivePlayerID: activePlayerID,
		Piece:          winnerPiece,
	}
}

func parseError(data string) tea.Msg {
	return protocol.ServerErrorMsg{
		Msg: data,
	}
}
