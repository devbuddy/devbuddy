package executor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandFalse(t *testing.T) {
	code, err := New("false").Run()

	require.Error(t, err)
	require.Equal(t, 1, code)
}

func TestCommandTrue(t *testing.T) {
	code, err := New("true").Run()

	require.NoError(t, err)
	require.Equal(t, 0, code)
}

func TestShellTrue(t *testing.T) {
	code, err := NewShell("true").Run()

	require.NoError(t, err)
	require.Equal(t, 0, code)
}

func TestShellFalse(t *testing.T) {
	code, err := NewShell("false").Run()

	require.Error(t, err)
	require.Equal(t, 1, code)
}

func TestShellCapture(t *testing.T) {
	output, code, err := NewShell("echo poipoi").Capture()

	require.NoError(t, err)
	require.Zero(t, code)
	require.Equal(t, "poipoi\n", output)
}

func TestShellCapturePWD(t *testing.T) {
	output, code, err := NewShell("echo $PWD").SetCwd("/bin").Capture()

	require.NoError(t, err)
	require.Zero(t, code)
	require.Equal(t, "/bin\n", output)
}

func TestCommandNotFound(t *testing.T) {
	code, err := New("never-ever-cmd").Run()

	require.Error(t, err)
	require.Equal(t, -1, code)
}
