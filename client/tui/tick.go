package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type FallTickMsg struct{}
type RowClearTickMsg struct{}

func fallTick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return FallTickMsg{}
	})
}

func clearRowTick() tea.Cmd {
	return tea.Tick(
		75*time.Millisecond,
		func(t time.Time) tea.Msg {
			return RowClearTickMsg{}
		},
	)
}
