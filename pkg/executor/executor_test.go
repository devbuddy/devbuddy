package executor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandFalse(t *testing.T) {
	exec := NewExecutor()
	result := exec.Run(New("false"))

	require.Error(t, result.Error)
	require.Equal(t, "command failed with exit code 1", result.Error.Error())
	require.NoError(t, result.LaunchError)
	require.Equal(t, 1, result.Code)
	require.Equal(t, "", result.Output)
}

func TestCommandTrue(t *testing.T) {
	exec := NewExecutor()
	result := exec.Run(New("true"))

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "", result.Output)
}

func TestShellTrue(t *testing.T) {
	exec := NewExecutor()
	result := exec.Run(NewShell("true"))

	require.NoError(t, result.Error)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "", result.Output)
}

func TestShellFalse(t *testing.T) {
	exec := NewExecutor()
	result := exec.Run(NewShell("false"))

	require.Error(t, result.Error)
	require.Equal(t, "command failed with exit code 1", result.Error.Error())
	require.NoError(t, result.LaunchError)
	require.Equal(t, 1, result.Code)
	require.Equal(t, "", result.Output)
}

func TestShellCapture(t *testing.T) {
	exec := NewExecutor()
	result := exec.Capture(NewShell("echo poipoi"))

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "poipoi\n", result.Output)
}

func TestShellCaptureAndTrim(t *testing.T) {
	exec := NewExecutor()
	result := exec.CaptureAndTrim(NewShell("echo poipoi"))

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "poipoi", result.Output)
}

func TestShellCapturePWD(t *testing.T) {
	exec := NewExecutor()
	cmd := NewShell("echo $PWD")
	cmd.Cwd = "/opt"
	result := exec.Capture(cmd)

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "/opt\n", result.Output)
}

func TestCommandNotFound(t *testing.T) {
	exec := NewExecutor()
	result := exec.Run(New("cmd-that-does-not-exist"))

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
	exec := NewExecutor()
	cmd := NewShell("echo $POIPOI")
	cmd.Env = []string{"POIPOI=something"}
	result := exec.Capture(cmd)

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "something\n", result.Output)
}

func TestAddEnvVar(t *testing.T) {
	exec := NewExecutor()
	result := exec.Capture(NewShell("echo ${V1}-${V2}").AddEnvVar("V1", "v1").AddEnvVar("V2", "v2"))

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "v1-v2\n", result.Output)
}

func TestPrefix(t *testing.T) {
	exec := NewExecutor()
	buf := &bytes.Buffer{}

	cmd := NewShell("echo \"line1\nline2\nline3\"")
	cmd.OutputWriter = buf
	cmd.OutputPrefix = "---"
	result := exec.Run(cmd)

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "---line1\n---line2\n---line3\n", buf.String())
}

func TestFilter(t *testing.T) {
	exec := NewExecutor()
	buf := &bytes.Buffer{}

	cmd := NewShell("echo \"line1\nline2\nline3\nline4\"")
	cmd.OutputWriter = buf
	cmd.AddOutputFilter("line2")
	cmd.AddOutputFilter("line4")
	result := exec.Run(cmd)

	require.NoError(t, result.Error)
	require.NoError(t, result.LaunchError)
	require.Equal(t, 0, result.Code)
	require.Equal(t, "line1\nline3\n", buf.String())
}
