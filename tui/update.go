package tui

import (
	tea "charm.land/bubbletea/v2"
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
	case tickMsg:
		m.game.MoveDown()
		return m, tick()

	case tea.KeyPressMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "a":
			m.game.MoveLeft()

		case "d":
			m.game.MoveRight()

		case "s":
			m.game.MoveDown()

		case "space":
			m.game.HardDrop()
		}

	}
	return m, nil
}
