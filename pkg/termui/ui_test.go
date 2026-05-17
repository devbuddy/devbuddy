package termui

import (
	"testing"

	baseui "github.com/devbuddy/devbuddy/pkg/ui"
	"github.com/stretchr/testify/require"
)

func TestUIRecordsEventsWhileRendering(t *testing.T) {
	buf, ui := NewTesting(false)

	ui.JumpProject("github.com/org/repo")

	require.Contains(t, buf.String(), "jumping to")
	require.Equal(t, []baseui.Event{
		{Kind: baseui.KindJumpProject, Text: "github.com/org/repo"},
	}, ui.Events())
}

func TestUIDoesNotRecordDisabledDebugEvents(t *testing.T) {
	_, ui := NewTesting(false)

	ui.Debug("hidden")

	require.Empty(t, ui.Events())
}

func TestUIRecordsTaskFields(t *testing.T) {
	_, ui := NewTesting(false)

	ui.TaskHeader("python", "3.9.0", "disabled")

	require.Equal(t, []baseui.Event{
		{
			Kind: baseui.KindTaskHeader,
			Text: "python",
			Fields: []baseui.Field{
				baseui.F("param", "3.9.0"),
				baseui.F("reason", "disabled"),
			},
		},
	}, ui.Events())
}
