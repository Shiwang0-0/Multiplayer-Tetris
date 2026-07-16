package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Shiwang0-0/multiplayertetris/game"
	"github.com/Shiwang0-0/multiplayertetris/tui"
)

func main() {

	// initalize a game

	game := game.NewGame()

	tuiModel := tui.NewModel(game)

	// initalizing the bubble tea program
	p := tea.NewProgram(tuiModel)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
