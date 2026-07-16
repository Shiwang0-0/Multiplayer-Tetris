package game

type Game struct {
	board       [Rows][Cols]Cell
	activePiece Piece
	gameOver    bool
}

func NewGame() *Game {
	g := &Game{
		board:    initializeBoard(),
		gameOver: false,
	}
	g.SpawnNewPiece()

	return g
}

func initializeBoard() [Rows][Cols]Cell {

	var board [Rows][Cols]Cell
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			switch {
			case i == 0 || i == Rows-1:
				board[i][j].Val = HorizontalBoundary

			case j == 0 || j == Cols-1:
				board[i][j].Val = VerticalBoundary

			default:
				board[i][j].Val = Empty
			}
		}
	}
	return board
}

func (g *Game) GetBoard() [Rows][Cols]Cell {
	return g.board // copy on assignment/return
}

func (g *Game) GetActivePiece() Piece {
	return g.activePiece
}

func (g *Game) IsGameOver() bool {
	return g.gameOver
}
