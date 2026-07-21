package game

func (g *Game) MoveLeft() {
	if g.CanMove(0, -1) {
		g.activePiece.AnchorY--
	}
}

func (g *Game) MoveRight() {
	if g.CanMove(0, 1) {
		g.activePiece.AnchorY++
	}
}

func (g *Game) MoveDown() bool {
	if g.CanMove(1, 0) {
		g.activePiece.AnchorX++
		return false
	}
	g.lockAndAdvance()
	return true
}

func (g *Game) lockAndAdvance() {
	g.LockPiece()
	g.clearRows = g.FindCompleteRows()

	if len(g.clearRows) > 0 {
		g.gameState = Clearing
		return
	}
	g.gameState = AwaitingTurn
}

func (g *Game) CanMove(dx, dy int) bool {
	for _, cell := range GetShapeCells(g.activePiece.Shape) {
		nextX := g.activePiece.AnchorX + cell.X + dx
		nextY := g.activePiece.AnchorY + cell.Y + dy

		if nextX <= 0 || nextX >= Rows-1 {
			return false
		}

		if nextY <= 0 || nextY >= Cols-1 {
			return false
		}

		// already a block present
		if g.board[nextX][nextY].Val == Block {
			return false
		}
	}

	return true
}

func (g *Game) HardDrop() {
	for g.CanMove(1, 0) {
		g.activePiece.AnchorX++
	}
	g.lockAndAdvance()
}
