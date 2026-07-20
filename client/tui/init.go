package tui

import tea "charm.land/bubbletea/v2"

func (m model) Init() tea.Cmd {
	return nil // the dropping of blocks starts when use is done waiting for more players to join the room and renders GameScreen
}
