package ui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUIRecordsEvents(t *testing.T) {
	ui := New()

	ui.Record(Event{Kind: KindJumpProject, Text: "org/repo"})

	require.Equal(t, []Event{{Kind: KindJumpProject, Text: "org/repo"}}, ui.Events())
}

func TestUIEventsReturnsCopy(t *testing.T) {
	ui := New()
	ui.Record(Event{Kind: KindWarning, Text: "careful"})

	events := ui.Events()
	events[0].Text = "mutated"

	require.Equal(t, []Event{{Kind: KindWarning, Text: "careful"}}, ui.Events())
}
