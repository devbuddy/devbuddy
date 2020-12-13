package expect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShellExpect_Bash(t *testing.T) {
	env := []string{
		"PS1=##\n",
		"TESTVAR=foobar",
	}

	ep := NewExpectWithEnv("bash", []string{"--noprofile", "--norc"}, env)
	err := ep.Start()
	require.NoError(t, err)

	shell := NewShellExpect(ep, "##\n")

	err = shell.Init()
	require.NoError(t, err)

	output, err := shell.Run("echo $TESTVAR")
	require.NoError(t, err)
	require.Equal(t, []string{"foobar"}, output)
}

func Test_ShellExpect_Zsh(t *testing.T) {
	t.SkipNow()

	env := []string{
		"PROMPT=##\n",
		"TESTVAR=foobar",
	}

	ep := NewExpectWithEnv("zsh", []string{"--no-globalrcs", "--no-rcs", "--no-zle", "--no-promptcr"}, env)
	err := ep.Start()
	require.NoError(t, err)
	ep.Debug = true

	shell := NewShellExpect(ep, "##\n")

	err = shell.Init()
	require.NoError(t, err)

	output, err := shell.Run("echo $TESTVAR")
	require.NoError(t, err)
	require.Equal(t, []string{"foobar"}, output)
}

const DockerImage = "devbuddy-test-env-linux"

func Test_ShellExpect_Docker_Bash(t *testing.T) {
	args := []string{
		"docker", "run", "-ti", "--rm",
		"-e", "PS1=##\n",
		"-e", "TESTVAR=foobar",
		"--entrypoint", "/bin/bash",
		DockerImage,
		"--noprofile", "--norc",
	}

	ep := NewExpect(args[0], args[1:]...)
	err := ep.Start()
	require.NoError(t, err)
	// ep.Debug = true

	shell := NewShellExpect(ep, "##\n")

	err = shell.Init()
	require.NoError(t, err)

	output, err := shell.Run("echo $BASH_VERSION")

	output, err = shell.Run("stty -echo") // disable echo inside the container
	require.NoError(t, err)
	require.Equal(t, []string{"stty -echo"}, output)

	output, err = shell.Run("echo $TESTVAR")
	require.NoError(t, err)
	require.Equal(t, []string{"foobar"}, output)
}

func Test_ShellExpect_Docker_Zsh(t *testing.T) {
	args := []string{
		"docker", "run", "-ti", "--rm",
		"-e", "PROMPT=##\n",
		"-e", "TESTVAR=foobar",
		"--entrypoint", "/bin/zsh",
		DockerImage,
		"--no-globalrcs", "--no-rcs", "--no-zle", "--no-promptcr",
	}

	ep := NewExpect(args[0], args[1:]...)
	err := ep.Start()
	require.NoError(t, err)
	ep.Debug = true

	shell := NewShellExpect(ep, "##\n")

	err = shell.Init()
	require.NoError(t, err)

	output, err := shell.Run("stty -echo") // disable echo inside the container
	require.NoError(t, err)
	require.Equal(t, []string{"stty -echo"}, output)

	output, err = shell.Run("echo $TESTVAR")
	require.NoError(t, err)
	require.Equal(t, []string{"foobar"}, output)
}
