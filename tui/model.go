package tui

import (
	"github.com/Shiwang0-0/multiplayertetris/game"
)

type model struct {
	game *game.Game
}

func NewModel(game *game.Game) *model {
	return &model{
		game: game,
	}
}
