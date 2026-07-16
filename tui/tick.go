package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
