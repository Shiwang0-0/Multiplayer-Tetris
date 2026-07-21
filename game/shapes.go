package game

type Point struct {
	X int
	Y int
}

func GetShape(p string) Shape {
	var shape Shape
	switch p {
	case "L":
		shape = L
	case "I":
		shape = I
	case "T":
		shape = T
	case "O":
		shape = O
	}
	return shape
}

// to prevent shape going out of the boundary
func GetShapeCells(shape Shape) []Point {
	switch shape {
	case L:
		return []Point{
			{0, 0},
			{1, 0},
			{2, 0},
			{2, 1},
		}

	case I:
		return []Point{
			{0, 0},
			{1, 0},
			{2, 0},
			{3, 0},
		}

	case T:
		return []Point{
			{0, 0},
			{0, 1},
			{0, 2},
			{1, 1},
		}

	case O:
		return []Point{
			{0, 0},
			{0, 1},
			{1, 0},
			{1, 1},
		}
	}

	return nil
}
