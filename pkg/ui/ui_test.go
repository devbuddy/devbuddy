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

func TestNewTestingInstallsFakePrompts(t *testing.T) {
	prompts, ui := NewTesting()

	prompts.SelectValue = "feature-a"
	got, err := ui.Prompts().Select(SelectRequest{Label: "Select"})

	require.NoError(t, err)
	require.Equal(t, "feature-a", got)
	require.Len(t, prompts.SelectRequests, 1)
}

func TestUIRecordsEventsWhileRendering(t *testing.T) {
	buf, ui := NewBufferedTesting(false)

	ui.JumpProject("github.com/org/repo")

	require.Contains(t, buf.String(), "jumping to")
	require.Equal(t, []Event{
		{Kind: KindJumpProject, Text: "github.com/org/repo"},
	}, ui.Events())
}

func TestUIDoesNotRecordDisabledDebugEvents(t *testing.T) {
	_, ui := NewBufferedTesting(false)

	ui.Debug("hidden")

	require.Empty(t, ui.Events())
}

func TestUIRecordsTaskFields(t *testing.T) {
	_, ui := NewBufferedTesting(false)

	ui.TaskHeader("python", "3.9.0", "disabled")

	require.Equal(t, []Event{
		{
			Kind: KindTaskHeader,
			Text: "python",
			Fields: []Field{
				F("param", "3.9.0"),
				F("reason", "disabled"),
			},
		},
	}, ui.Events())
}
