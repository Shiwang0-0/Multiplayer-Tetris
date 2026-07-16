package game

const (
	HorizontalBoundary rune = '-'
	VerticalBoundary   rune = '|'
	Empty              rune = ' '
	Block              rune = '█'
	Rows               int  = 22
	Cols               int  = 12
)

type Shape byte

const (
	L Shape = iota
	I
	T
	O
)

type Cell struct {
	Color ColorID
	Val   rune
}

type Piece struct {
	Shape   Shape
	ColorID ColorID
	AnchorX int
	AnchorY int
}
