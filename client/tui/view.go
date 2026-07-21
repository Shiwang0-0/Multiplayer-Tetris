package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Shiwang0-0/multiplayertetris/game"
)

const (
	HomeScreen Screen = iota
	CreateRoomScreen
	JoinRoomScreen
	GameScreen
	WaitingScreen
	StartScreen
	VotingScreen
	SpectateScreen
)

func (m model) View() tea.View {

	switch m.screen {

	case HomeScreen:
		return tea.NewView(
			"Multiplayer Tetris\n\n" +
				"[C] Create Room\n" +
				"[J] Join Room\n",
		)

	case CreateRoomScreen:
		return tea.NewView(
			"Create Room\n\n" +
				"Room ID: " + m.roomID + "\n\n" +
				"Enter: Create | Esc: Back",
		)

	case JoinRoomScreen:
		return tea.NewView(
			"Join Room\n\n" +
				"Room ID: " + m.roomID + "\n\n" +
				"Enter: Join | Esc: Back",
		)

	case WaitingScreen:
		return tea.NewView(
			"Waiting for players...\n\n" +
				"Room ID: " + m.roomID + "\n\n" +
				"Waiting for opponent...\n\n" +
				"Press Enter to start...\n\n" +
				"Press Esc to leave",
		)
	case SpectateScreen:
		opp, ok := m.opponents[m.activePlayerID]
		if !ok {
			return tea.NewView("Waiting for the active player's board...\n")
		}
		header := fmt.Sprintf("Spectating Player %d's turn\n\n", m.activePlayerID)
		return tea.NewView(header + m.renderBoardFor(opp))

	case VotingScreen:
		var boardView string
		if opp, ok := m.opponents[m.activePlayerID]; ok {
			boardView = m.renderBoardFor(opp)
		} else if m.activePlayerID == m.myID {
			boardView = m.renderBoard()
		}

		header := fmt.Sprintf(
			"Voting for Player %d's next piece (%ds left)\n\n",
			m.activePlayerID, m.voteSecondsLeft,
		)

		prompt := "Your vote: "
		switch {
		case m.myVote != "":
			prompt += m.myVote
		case m.myID == m.activePlayerID:
			prompt += "(you can't vote for your own piece)"
		default:
			prompt += "[1] I  [2] O  [3] T  [4] L"
		}

		return tea.NewView(header + boardView + "\n\n" + prompt)

	case GameScreen:
		return tea.NewView(m.renderBoard())
	}
	return tea.NewView("")
}

func (m model) renderBoard() string {
	return m.renderBoardFor(m.game)
}

func (m model) renderBoardFor(g *game.Game) string {
	var board strings.Builder
	tempBoard := g.GetBoard()
	m.renderActivePieceFor(g, &tempBoard)

	for i := 0; i < game.Rows; i++ {
		for j := 0; j < game.Cols; j++ {
			flashStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

			if j == 0 || j == game.Cols-1 {
				board.WriteString("|")
				continue
			}

			if g.GetGameState() == game.Clearing && g.IsClearRow(i) {
				if g.GetClearRowTick()%2 == 0 {
					board.WriteString(flashStyle.Render("██"))
				} else {
					board.WriteString("  ")
				}
				continue
			}

			cell := tempBoard[i][j]
			switch cell.Val {
			case game.Block:
				style := lipgloss.NewStyle().Foreground(getLipglossColor(cell.Color))
				board.WriteString(style.Render("██"))
			case game.HorizontalBoundary:
				board.WriteString("--")
			case game.VerticalBoundary:
				board.WriteString("|")
			default:
				board.WriteString("  ")
			}
		}
		board.WriteRune('\n')
	}

	if g.IsGameOver() {
		board.WriteString("\n")
		board.WriteString(m.renderGameOver())
	}

	return board.String()
}

func (m model) renderActivePieceFor(g *game.Game, board *[game.Rows][game.Cols]game.Cell) {
	activePiece := g.GetActivePiece()
	for _, cell := range game.GetShapeCells(activePiece.Shape) {
		x := activePiece.AnchorX + cell.X
		y := activePiece.AnchorY + cell.Y
		board[x][y] = game.Cell{Val: game.Block, Color: activePiece.ColorID}
	}
}

func (m model) renderGameOver() string {
	var board strings.Builder

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9")).
		Width(game.Cols * 2).
		Align(lipgloss.Center)

	board.WriteString(style.Render("GAME OVER"))
	board.WriteRune('\n')
	board.WriteString("Press q to quit")

	return board.String()
}
