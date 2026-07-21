package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type FallTickMsg struct{}
type RowClearTickMsg struct{}
type VoteCountdownTickMsg struct{}

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

// how much time is left for this current voting session
func voteCountdownTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(time.Time) tea.Msg { // ticks every second, and stops after 10th second
		return VoteCountdownTickMsg{}
	})
}
