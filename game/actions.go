package game

func (g *Game) LockPiece() {
	for _, cell := range GetShapeCells(g.activePiece.Shape) {
		x := g.activePiece.AnchorX + cell.X
		y := g.activePiece.AnchorY + cell.Y

		g.board[x][y] = Cell{
			Val:   Block,
			Color: g.activePiece.ColorID,
		}
	}
}

func (g *Game) FindCompleteRows() []int {
	var lines []int

	for i := Rows - 2; i > 0; i-- {
		complete := true

		for j := 1; j < Cols-1; j++ {
			if g.board[i][j].Val != Block {
				complete = false
				break
			}
		}

		if complete {
			lines = append(lines, i)
		}
	}

	return lines
}

func (g *Game) ClearCompleteRows() {
	clearSet := make(map[int]bool)

	for _, row := range g.clearRows {
		clearSet[row] = true
	}

	for i := Rows - 2; i > 0; i-- {
		if clearSet[i] {
			continue
		}

		shift := 0

		for _, clearRow := range g.clearRows {
			if clearRow > i {
				shift++
			}
		}

		newRow := i + shift

		for j := 1; j < Cols-1; j++ {
			g.board[newRow][j] = g.board[i][j]
		}
	}

	// Clear top rows.
	for i := 1; i <= len(g.clearRows); i++ {
		for j := 1; j < Cols-1; j++ {
			g.board[i][j] = Cell{
				Val: Empty,
			}
		}
	}

	g.clearRows = nil
}

func (g *Game) UpdateClearAnimation() {
	g.clearRowTick++

	if g.clearRowTick >= 3 {
		g.ClearCompleteRows()
		g.clearRows = nil
		g.clearRowTick = 0

		g.SpawnNewPiece()
		g.gameState = Playing
	}
}

func (g *Game) SpawnNewPiece() {
	piece := Piece{
		Shape:   GetShape(),
		ColorID: GetColor(),
		AnchorX: 1,
		AnchorY: Cols / 2,
	}
	g.activePiece = piece
}
