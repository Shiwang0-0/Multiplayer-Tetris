package game

import "math/rand"

type Point struct {
	X int
	Y int
}

func GetShape() Shape {
	shapes := []Shape{L, I, T, O}

	return shapes[rand.Intn(len(shapes))]
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
