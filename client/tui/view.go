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
	EliminatedScreen
	MatchOverScreen
)

func (m model) View() tea.View {

	switch m.screen {

	case HomeScreen:
		body := titleStyle.Render("♦ Multiplayer Tetris ♦") + "\n\n" +
			"[C] Create Room\n" +
			"[J] Join Room\n" +
			helpStyle.Render("q: quit")
		if m.lastError != "" {
			body += "\n" + errorStyle.Render(m.lastError)
		}
		return m.center(boxStyle.Render(body))

	case CreateRoomScreen:
		body := titleStyle.Render("Create Room") + "\n\n" +
			"Room ID: " + m.roomID + "\n\n" +
			helpStyle.Render("enter: create   esc: back")
		if m.lastError != "" {
			body += "\n" + errorStyle.Render(m.lastError)
		}
		return m.center(boxStyle.Render(body))

	case JoinRoomScreen:
		body := titleStyle.Render("Join Room") + "\n\n" +
			"Room ID: " + m.roomID + "\n\n" +
			helpStyle.Render("enter: join   esc: back")
		if m.lastError != "" {
			body += "\n" + errorStyle.Render(m.lastError)
		}
		return m.center(boxStyle.Render(body))

	case WaitingScreen:
		body := titleStyle.Render("Waiting for players...") + "\n\n" +
			"Room ID: " + m.roomID + "\n\n" +
			subtleStyle.Render("Waiting for opponent...")

		if m.isCreator {
			body += "\n\n" + helpStyle.Render("enter: start   esc: leave")
		} else {
			body += "\n\n" + helpStyle.Render("esc: leave")
		}
		if m.lastError != "" {
			body += "\n" + errorStyle.Render(m.lastError)
		}
		return m.center(boxStyle.Render(body))

	case SpectateScreen:
		if m.eliminationNotice != "" {
			content := lipgloss.JoinVertical(lipgloss.Center,
				titleStyle.Render(m.eliminationNotice),
				subtleStyle.Render("waiting for the next turn..."),
			)
			return m.center(content)
		}
		opp, ok := m.opponents[m.activePlayerID]
		if !ok {
			return m.center(boxStyle.Render("Waiting for the active player's board..."))
		}
		header := titleStyle.Render(fmt.Sprintf("Spectating Player %d", m.activePlayerID))
		return m.center(header + "\n" + boardBorderStyle.Render(m.renderBoardFor(opp)))

	case VotingScreen:

		var boardView string
		if opp, ok := m.opponents[m.activePlayerID]; ok {
			boardView = boardBorderStyle.Render(m.renderBoardFor(opp))
		} else if m.activePlayerID == m.myID {
			boardView = boardBorderStyle.Render(m.renderBoard())
		}

		header := titleStyle.Render(
			fmt.Sprintf("Voting for Player %d's next piece  (%ds left)",
				m.activePlayerID, m.voteSecondsLeft),
		)

		var prompt string
		switch {
		case m.myVote != "":
			prompt = "Your vote: " + activeVoteStyle.Render(m.myVote)
		case m.myID == m.activePlayerID:
			prompt = subtleStyle.Render("(you can't vote for your own piece)")
		default:
			prompt = "Your vote: [1] I   [2] O   [3] T   [4] L"
		}

		content := lipgloss.JoinVertical(lipgloss.Center,
			header, boardView, "", prompt,
		)
		return m.center(content)

	case EliminatedScreen:
		content := lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("You're out!"),
			boardBorderStyle.Render(m.renderBoard()), // still shows your final board + the GAME OVER stamp
			helpStyle.Render("s: spectate   q: quit"),
		)
		return m.center(content)

	case GameScreen:
		header := titleStyle.Render("Your Board")
		content := lipgloss.JoinVertical(lipgloss.Center,
			header,
			boardBorderStyle.Render(m.renderBoard()),
			helpStyle.Render("a: left   d: right   s: soft drop   space: hard drop"),
		)
		return m.center(content)

	case MatchOverScreen:
		var msg string
		if m.winnerID == m.myID {
			msg = "🏆 You won!"
		} else {
			msg = fmt.Sprintf("Player %d wins!", m.winnerID)
		}
		content := lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("Match Over"),
			msg,
			helpStyle.Render("q: quit"),
		)
		return m.center(content)
	}
	return m.center("")
}

func (m model) renderBoard() string {
	return m.renderBoardFor(m.game)
}

func (m model) renderBoardFor(g *game.Game) string {
	var board strings.Builder
	tempBoard := g.GetBoard()
	if g.GetGameState() == game.Playing {
		m.renderActivePieceFor(g, &tempBoard)
	}

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
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9")).
		Width(game.Cols * 2).
		Align(lipgloss.Center)

	return style.Render("GAME OVER") + "\n" + subtleStyle.Render("Press q to quit")
}

func (m model) center(content string) tea.View {
	var v tea.View
	if m.width == 0 || m.height == 0 {
		v = tea.NewView(content)
	} else {
		v = tea.NewView(
			lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content),
		)
	}
	v.AltScreen = true
	return v
}
