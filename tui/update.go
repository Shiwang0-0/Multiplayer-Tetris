package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Shiwang0-0/multiplayertetris/game"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.game.IsGameOver() {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
		}

		return m, nil
	}

	switch msg := msg.(type) {
	case FallTickMsg:
		if m.game.GetGameState() == game.Playing {
			m.game.MoveDown()
		}

		if m.game.GetGameState() == game.Clearing {
			return m, clearRowTick()
		}

		return m, fallTick()

	case RowClearTickMsg:
		if m.game.GetGameState() == game.Clearing {
			m.game.UpdateClearAnimation()
			return m, clearRowTick()
		}

		return m, fallTick()

	case tea.KeyPressMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "a":
			if m.game.GetGameState() == game.Playing {
				m.game.MoveLeft()
			}

		case "d":
			if m.game.GetGameState() == game.Playing {
				m.game.MoveRight()
			}

		case "s":
			if m.game.GetGameState() == game.Playing {
				m.game.MoveDown()
			}

		case "space":
			if m.game.GetGameState() == game.Playing {
				m.game.HardDrop()
			}
		}
	}
	return m, nil
}
