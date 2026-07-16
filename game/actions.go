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
func (g *Game) SpawnNewPiece() {
	piece := Piece{
		Shape:   GetShape(),
		ColorID: GetColor(),
		AnchorX: 1,
		AnchorY: Cols / 2,
	}
	g.activePiece = piece
}
