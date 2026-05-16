package harness

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// OutputContains asserts that every substring appears in the joined output.
func OutputContains(t *testing.T, lines []string, subStrings ...string) {
	t.Helper()
	text := StripAnsi(strings.Join(lines, "\n"))
	for _, s := range subStrings {
		require.Contains(t, text, s)
	}
}

// OutputNotContains asserts that no substring appears in the joined output.
func OutputNotContains(t *testing.T, lines []string, subStrings ...string) {
	t.Helper()
	text := StripAnsi(strings.Join(lines, "\n"))
	for _, s := range subStrings {
		require.NotContains(t, text, s)
	}
}

// OutputEqual asserts that the lines match expectedLines exactly.
func OutputEqual(t *testing.T, lines []string, expectedLines ...string) {
	t.Helper()
	require.Equal(t, expectedLines, lines)
}
