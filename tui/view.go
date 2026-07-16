package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Shiwang0-0/multiplayertetris/game"
)

func (m model) View() tea.View {

	return tea.NewView(m.renderBoard())
}

func (m model) renderBoard() string {

	var board strings.Builder // for temporary state while the block is coming down

	// Temporary board for rendering the active piece.
	tempBoard := m.game.GetBoard()

	m.renderActivePiece(&tempBoard) // renders the falling down peice (doesnt change the board)

	for i := 0; i < game.Rows; i++ {
		for j := 0; j < game.Cols; j++ {
			cell := tempBoard[i][j]
			switch cell.Val {

			case game.Block:
				style := lipgloss.NewStyle().
					Foreground(getLipglossColor(cell.Color))

				board.WriteString(style.Render("██"))

			case game.HorizontalBoundary:
				board.WriteString("--")

			case game.VerticalBoundary:
				board.WriteString("| ")

			default:
				board.WriteString("  ")
			}
		}

		board.WriteRune('\n')
	}

	if m.game.IsGameOver() {
		board.WriteString("\n")
		board.WriteString(m.renderGameOver())
	}

	return board.String()
}

func (m model) renderActivePiece(board *[game.Rows][game.Cols]game.Cell) {

	activePiece := m.game.GetActivePiece()

	for _, cell := range game.GetShapeCells(activePiece.Shape) {
		x := activePiece.AnchorX + cell.X
		y := activePiece.AnchorY + cell.Y

		board[x][y] = game.Cell{
			Val:   game.Block,
			Color: activePiece.ColorID,
		}
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
