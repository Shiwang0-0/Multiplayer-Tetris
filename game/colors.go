package game

import (
	"math/rand"
)

type ColorID byte

const (
	Red ColorID = iota
	Green
	Yellow
	Blue
	Magenta
	Cyan
)

func GetColor() ColorID {
	colors := []ColorID{Red, Green, Yellow, Blue, Magenta, Cyan}

	return colors[rand.Intn(len(colors))]
}
