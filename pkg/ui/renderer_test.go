package ui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlainRendererDoesNotEmitANSI(t *testing.T) {
	got := PlainRenderer{}.Render(Event{Kind: KindJumpProject, Text: "org/repo"})

	require.Equal(t, "🐼  jumping to org/repo\n", got)
}

func TestTerminalRendererStylesOutput(t *testing.T) {
	got := TerminalRenderer{}.Render(Event{Kind: KindJumpProject, Text: "org/repo"})

	require.Contains(t, got, "\x1b[")
	require.Contains(t, got, "jumping to")
	require.Contains(t, got, "org/repo")
}

func TestRendererFormatsFeatureList(t *testing.T) {
	event := Event{
		Kind: KindHookActivated,
		Fields: []Field{
			F("python", "3.9.0"),
			F("env", ""),
		},
	}

	require.Equal(t, "🐼  activated: python[3.9.0], env\n", PlainRenderer{}.Render(event))
}
