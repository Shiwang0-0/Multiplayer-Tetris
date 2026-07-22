package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/Shiwang0-0/multiplayertetris/protocol"
)

type FallTickMsg struct{}
type RowClearTickMsg struct{}
type VoteCountdownTickMsg struct{}

func fallTick() tea.Cmd {
	return tea.Tick(protocol.FallDuration, func(t time.Time) tea.Msg {
		return FallTickMsg{}
	})
}

func clearRowTick() tea.Cmd {
	return tea.Tick(
		protocol.ClearRowDuration,
		func(t time.Time) tea.Msg {
			return RowClearTickMsg{}
		},
	)
}

// how much time is left for this current voting session
func voteCountdownTick() tea.Cmd {
	return tea.Tick(protocol.VoteCountdownDuration, func(time.Time) tea.Msg { // ticks every second, and stops after 10th second
		return VoteCountdownTickMsg{}
	})
}
