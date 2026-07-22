package tui

import (
	"fmt"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"github.com/Shiwang0-0/multiplayertetris/game"
	"github.com/Shiwang0-0/multiplayertetris/protocol"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case FallTickMsg:
		if m.game.GetGameState() == game.Playing {
			locked := m.game.MoveDown()
			if locked {
				return m, tea.Batch(fallTick(), m.sendCommand("DOWN"), m.sendCommand("LOCKED"))
			}
			return m, tea.Batch(fallTick(), m.sendCommand("DOWN"))
		}
		if m.game.GetGameState() == game.Clearing {
			return m, clearRowTick()
		}
		return m, nil

	case RowClearTickMsg:
		if m.game.GetGameState() == game.Clearing {
			m.game.UpdateClearAnimation()
			if m.game.GetGameState() != game.Clearing {
				// animation just finished this tick
				return m, tea.Batch(fallTick(), m.sendCommand("CLEARED"))
			}
			return m, clearRowTick()
		}
		return m, fallTick()

	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)

	case protocol.ClientIDMsg:
		m.myID = msg.ID
		return m, nil

	case protocol.JoinedMsg:
		m.screen = WaitingScreen
		return m, nil

	case protocol.VotingStartMsg:
		if m.eliminated {
			return m, nil
		}
		m.screen = VotingScreen
		m.activePlayerID = msg.ActivePlayerID
		m.myVote = ""
		m.voteSecondsLeft = msg.DeadlineSeconds
		return m, voteCountdownTick()

	case protocol.TurnStartMsg:
		m.eliminationNotice = ""
		m.activePlayerID = msg.ActivePlayerID

		if msg.ActivePlayerID == m.myID {
			m.game.SpawnNewPiece(msg.Piece)
			if m.game.IsGameOver() { // represents game of the current player, not t
				m.eliminated = true
				m.screen = EliminatedScreen
				return m, m.sendCommand("OUT")
			}
			m.screen = GameScreen
			return m, fallTick()
		}

		if m.opponents == nil {
			m.opponents = map[int]*game.Game{}
		}
		opp, ok := m.opponents[msg.ActivePlayerID]
		if !ok {
			opp = game.NewGame()
			m.opponents[msg.ActivePlayerID] = opp
		}
		opp.SpawnNewPiece(msg.Piece) // keeps the spectated board's piece in sync with the real active board
		if m.screen != EliminatedScreen {
			m.screen = SpectateScreen
		}
		return m, nil

	case VoteCountdownTickMsg:
		if m.screen != VotingScreen {
			return m, nil // voting already ended (TURN_START arrived)
		}
		if m.voteSecondsLeft > 0 {
			m.voteSecondsLeft--
			return m, voteCountdownTick() // schedule next second
		}

		return m, nil // countdown finished

	case protocol.ServerErrorMsg:
		m.lastError = msg.Msg
		return m, nil

	case protocol.OpponentMoveMsg:
		if msg.SenderID == m.myID {
			return m, nil // ignore, in case the server ever echoes our own move back
		}
		if m.opponents == nil {
			m.opponents = map[int]*game.Game{}
		}
		opp, ok := m.opponents[msg.SenderID] // this state is replicated here on the client side (state is not bought here, only similar movements are being played)
		if !ok {
			opp = game.NewGame()
			m.opponents[msg.SenderID] = opp
		}
		switch msg.Move {
		case "LEFT":
			opp.MoveLeft()
		case "RIGHT":
			opp.MoveRight()
		case "DOWN":
			opp.MoveDown()
		case "DROP":
			opp.HardDrop()
		case "CLEARED":
			opp.ResolveClear()
		}
		return m, nil

	case protocol.PlayerOutMsg:
		if msg.PlayerID != m.myID {
			m.eliminationNotice = fmt.Sprintf("Player %d is out!", msg.PlayerID)
		}
		return m, nil

	case protocol.MatchOverMsg:
		m.matchOver = true
		m.winnerID = msg.WinnerID
		m.screen = MatchOverScreen
		return m, nil

	case protocol.DisconnectedMsg:
		m.connected = false // add `connected bool` to model if not present, to show a banner in View()
		return m, nil
	}
	return m, nil
}

func (m model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	}

	switch m.screen {
	case HomeScreen:
		switch msg.String() {
		case "c":
			m.screen = CreateRoomScreen
		case "j":
			m.screen = JoinRoomScreen
		}
		return m, nil

	case CreateRoomScreen, JoinRoomScreen:
		return m.handleRoomInput(msg)

	case WaitingScreen:
		if m.lastError != "" && msg.String() == "esc" {
			m.screen = HomeScreen
			m.lastError = ""
			return m, nil
		}
		if !m.isCreator {
			return m, nil
		}
		switch msg.String() {
		case "enter":
			return m, m.sendCommand("START_MATCH")
		}
		return m, nil

	case VotingScreen:
		if m.myID == m.activePlayerID {
			return m, nil // active player just waits
		}
		var piece string
		switch msg.String() {
		case "1":
			piece = "I"
		case "2":
			piece = "O"
		case "3":
			piece = "T"
		case "4":
			piece = "L"
		}
		if piece != "" && m.myVote == "" {
			m.myVote = piece
			return m, m.sendCommand("VOTE " + piece)
		}
		return m, nil

	case EliminatedScreen:
		if msg.String() == "s" {
			m.screen = SpectateScreen
		}
		return m, nil
	}

	switch msg.String() {
	case "a":
		if m.game.GetGameState() == game.Playing {
			m.game.MoveLeft()
			return m, m.sendCommand("LEFT")
		}
	case "d":
		if m.game.GetGameState() == game.Playing {
			m.game.MoveRight()
			return m, m.sendCommand("RIGHT")
		}
	case "s":
		if m.game.GetGameState() == game.Playing {
			if m.game.MoveDown() {
				return m, tea.Batch(m.sendCommand("DOWN"), m.sendCommand("LOCKED"))
			}
			return m, m.sendCommand("DOWN")
		}
	case "space":
		if m.game.GetGameState() == game.Playing {
			m.game.HardDrop()
			return m, tea.Batch(m.sendCommand("DROP"), m.sendCommand("LOCKED"))
		}
	}
	return m, nil
}

// root based keystokes updates
func (m *model) handleRoomInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {

	case "esc":
		m.screen = HomeScreen
		m.roomID = ""
		return m, nil

	case "backspace":
		if len(m.roomID) > 0 {
			m.roomID = m.roomID[:len(m.roomID)-1]
		}
		return m, nil

	case "enter":
		_, err := strconv.Atoi(m.roomID)
		if err != nil {
			return m, nil
		}
		m.lastError = ""
		if m.screen == CreateRoomScreen {
			m.isCreator = true
			return m, m.sendCommand("CREATE_ROOM " + m.roomID)
		}
		// m.screen = WaitingScreen
		return m, m.sendCommand("JOIN_ROOM " + m.roomID)

	default:
		key := msg.String()

		if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
			m.roomID += key
		}
	}

	return m, nil
}

// sendCommand is a tea.Cmd factory: it writes a command to the server
// and runs as a side effect returned from Update(), rather than blocking Update() itself with a direct write.
func (m model) sendCommand(cmd string) tea.Cmd {
	conn := m.conn
	return func() tea.Msg {
		if _, err := fmt.Fprintln(conn, cmd); err != nil {
			return protocol.DisconnectedMsg{Err: err}
		}
		return nil
	}
}
