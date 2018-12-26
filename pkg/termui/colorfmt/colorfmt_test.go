package colorfmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func runIt(enableColor bool, text string, a ...interface{}) string {
	buf := new(bytes.Buffer)
	New(buf, enableColor).Printf(text, a...)
	return buf.String()
}

func TestTags(t *testing.T) {
	type testRun struct {
		text     string
		expected string
	}

	tests := []testRun{
		testRun{"__", "__"},

		testRun{"_{yellow}_", "_\x1b[0;33m_"},
		testRun{"_{yellow+bBuih:black+h}_", "_\x1b[0;1;5;4;7;93m\x1b[100m_"},

		testRun{"_{notAColor}_", "__"},
		testRun{"_{with space}_", "__"},

		testRun{"_{link}_", "_\x1b[0;1;4;92m_"},
	}

	for _, run := range tests {
		require.Equal(t, run.expected, runIt(true, run.text))
		require.Equal(t, "__", runIt(false, run.text))
	}
}

func TestVariadic(t *testing.T) {
	require.Equal(t,
		"_\x1b[0;33mXX\x1b[0;32mYY_",
		runIt(true, "_{yellow}%s{green}%s_", "XX", "YY"),
	)
}

func TestEscapeBracket(t *testing.T) {
	require.Equal(t, "_{yellow}}_", runIt(true, "_{{yellow}}_"))
}

func TestMultipleTags(t *testing.T) {
	require.Equal(t, "_\x1b[0;33m_\x1b[0;32m_", runIt(true, "_{yellow}_{green}_"))
}

func TestReset(t *testing.T) {
	require.Equal(t, "_\x1b[0;33mXX\x1b[0mYY_", runIt(true, "_{yellow}XX{reset}YY_"))
}
