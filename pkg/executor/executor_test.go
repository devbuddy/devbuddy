package executor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandFalse(t *testing.T) {
	result := New("false").Run()

	require.Error(t, result.Error)
	require.Equal(t, "command failed with exit code 1", result.Error.Error())
	require.NoError(t, result.LaunchError)
	require.Equal(t, 1, result.Code)
	require.Equal(t, "", result.Output)
}

func TestCommandTrue(t *testing.T) {
	result := New("true").Run()

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "", result.Output)
}

func TestShellTrue(t *testing.T) {
	result := NewShell("true").Run()

	require.NoError(t, result.Error)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "", result.Output)
}

func TestShellFalse(t *testing.T) {
	result := NewShell("false").Run()

	require.Error(t, result.Error)
	require.Equal(t, "command failed with exit code 1", result.Error.Error())
	require.NoError(t, result.LaunchError)
	require.Equal(t, 1, result.Code)
	require.Equal(t, "", result.Output)
}

func TestShellCapture(t *testing.T) {
	result := NewShell("echo poipoi").Capture()

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "poipoi\n", result.Output)
}

func TestShellCaptureAndTrim(t *testing.T) {
	result := NewShell("echo poipoi").CaptureAndTrim()

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "poipoi", result.Output)
}

func TestShellCapturePWD(t *testing.T) {
	result := NewShell("echo $PWD").SetCwd("/opt").Capture()

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "/opt\n", result.Output)
}

func TestCommandNotFound(t *testing.T) {
	result := New("cmd-that-does-not-exist").Run()

	require.Error(t, result.Error)
	require.Equal(t,
		"command failed with: exec: \"cmd-that-does-not-exist\": executable file not found in $PATH",
		result.Error.Error())
	require.Error(t, result.LaunchError)
	require.Equal(t,
		"command failed with: exec: \"cmd-that-does-not-exist\": executable file not found in $PATH",
		result.LaunchError.Error())
	require.Equal(t, 0, result.Code)
	require.Equal(t, "", result.Output)
}

func TestSetEnv(t *testing.T) {
	result := NewShell("echo $POIPOI").SetEnv([]string{"POIPOI=something"}).Capture()

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "something\n", result.Output)
}

func TestSetEnvVar(t *testing.T) {
	result := NewShell("echo ${V1}-${V2}").SetEnvVar("V1", "v1").SetEnvVar("V2", "v2").Capture()

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "v1-v2\n", result.Output)
}

func TestPrefix(t *testing.T) {
	buf := &bytes.Buffer{}

	executor := NewShell("echo \"line1\nline2\nline3\"")
	impl := executor.(*executorImpl)
	impl.outputWriter = buf
	impl.SetOutputPrefix("---")
	result := impl.Run()

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "---line1\n---line2\n---line3\n", buf.String())
}

func TestFilter(t *testing.T) {
	buf := &bytes.Buffer{}

	executor := NewShell("echo \"line1\nline2\nline3\nline4\"")
	impl := executor.(*executorImpl)
	impl.outputWriter = buf
	impl.AddOutputFilter("line2")
	impl.AddOutputFilter("line4")
	result := impl.Run()

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "line1\nline3\n", buf.String())
}
