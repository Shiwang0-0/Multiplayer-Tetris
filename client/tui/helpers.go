package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/Shiwang0-0/multiplayertetris/game"
)

func getLipglossColor(colorID game.ColorID) color.Color {
	switch colorID {
	case game.Red:
		return lipgloss.Color("1")

	case game.Green:
		return lipgloss.Color("2")

	case game.Yellow:
		return lipgloss.Color("3")

	case game.Blue:
		return lipgloss.Color("4")

	case game.Magenta:
		return lipgloss.Color("5")

	case game.Cyan:
		return lipgloss.Color("6")
	}

	return lipgloss.Color("15")
}
