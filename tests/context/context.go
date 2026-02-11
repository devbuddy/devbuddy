package context

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/devbuddy/expect"
	"github.com/stretchr/testify/require"
)

type TestContext struct {
	expect *expect.ExpectProcess
	shell  interface {
		Run(string) ([]string, error)
	}
	debug bool
}

func New(config Config) (*TestContext, error) {
	var (
		shellPath string
		args      []string
	)

	// Maybe: export PAGER=cat; export SHELL=/bin/zsh

	switch config.ShellName {
	case "bash":
		shellPath = "/bin/bash"
		args = append(args, "--noprofile", "--norc")
	case "zsh":
		shellPath = "/bin/zsh"
		args = append(args, "--no-globalrcs", "--no-rcs", "--no-zle", "--no-promptcr")
	default:
		return nil, fmt.Errorf("unknown shell: %s", config.ShellName)
	}

	dockerExec := "docker"
	cmd := exec.Command("docker", "-v")
	if cmd.Run() != nil {
		dockerExec = "podman"
	}

	dockerCommand := []string{
		dockerExec, "run",
		"-ti",
		"-v", config.BinaryPath + ":/usr/local/bin/bud",
		"-e", "PROMPT=##\n",
		"-e", "PS1=##\n",
		"-e", "IN_DOCKER=yes",
		"--rm",
		"--entrypoint", shellPath,
		config.DockerImage,
	}
	dockerCommand = append(dockerCommand, args...)

	e := expect.NewExpect(dockerCommand[0], dockerCommand[1:]...)
	err := e.Start()
	if err != nil {
		return nil, fmt.Errorf("creating expect object: %w", err)
	}

	c := expect.NewShellExpect(e, "##\n")

	err = c.Init() // expect the initial prompt
	if err != nil {
		return nil, fmt.Errorf("running initialization shell command: %w", err)
	}

	tc := &TestContext{
		expect: e,
		shell:  c,
		// t:      t,
	}
	tc.debugLine("Expect command: %q", dockerCommand)

	_, err = tc.run("stty -echo")
	if err != nil {
		return nil, fmt.Errorf("disabling echo mode: %w", err)
	}

	output, err := tc.run("echo $IN_DOCKER")
	if err != nil {
		return nil, err
	}
	if len(output) != 1 || output[0] != "yes" {
		return nil, errors.New("not running in docker, IN_DOCKER var not found")
	}

	return tc, nil
}

func (c *TestContext) Verbose() {
	c.debug = true
	c.expect.Debug = false
}

func (c *TestContext) Debug() {
	c.debug = true
	c.expect.Debug = true
}

func (c *TestContext) Close() error {
	c.debugLine("Stopping docker container")
	err := c.expect.Stop()
	if err != nil {
		c.debugLine("ERROR when stopping docker: %s", err)
	}
	return err
}

func (c *TestContext) Run(t *testing.T, cmd string, optFns ...runOptionsFn) []string {
	t.Helper()
	lines, err := c.run(cmd, optFns...)
	require.NoError(t, err, "running command: %q", cmd)
	return lines
}

func (c *TestContext) run(cmd string, optFns ...runOptionsFn) ([]string, error) {
	opt := buildRunOptions(optFns)

	c.debugLine("Running command %q", cmd)
	c.debugLine("Options: %+v", opt)

	var lines []string
	var err error

	done := make(chan bool)
	go func() {
		lines, err = c.shell.Run(cmd)
		close(done)
	}()

	select {
	case <-done:
		c.debugLine("command completed")
	case <-time.After(opt.timeout):
		excerpt := "no output yet"
		if len(lines) > 0 {
			excerpt = strings.Join(lines[:min(10, len(lines))], "\n")
		}
		return nil, fmt.Errorf("timed out after %s running %q. output so far:\n%s", opt.timeout, cmd, excerpt)
	}

	codeLines, err := c.shell.Run("echo $?")
	if err != nil {
		return nil, fmt.Errorf("getting exit code: %w", err)
	}
	if len(codeLines) == 0 {
		return nil, fmt.Errorf("unexpected output when getting exit code: no output")
	}

	exitCode, err := strconv.Atoi(codeLines[0])
	if err != nil {
		return nil, fmt.Errorf("unexpected exit code %s: %w", codeLines[0], err)
	}
	if exitCode != opt.exitCode {
		excerpt := "no output"
		if len(lines) > 0 {
			excerpt = strings.Join(lines[:min(10, len(lines))], "\n")
		}
		return nil, fmt.Errorf("unexpected exit code %d (expected %d) running %q. output:\n%s", exitCode, opt.exitCode, cmd, excerpt)
	}

	return StripAnsiSlice(lines), nil
}

func (c *TestContext) Write(t *testing.T, path, content string) {
	t.Helper()
	b64content := base64.StdEncoding.EncodeToString([]byte(content))
	_, err := c.shell.Run(fmt.Sprintf("echo %s | base64 --decode > %q", b64content, path))
	require.NoError(t, err)
}

func (c *TestContext) WriteLines(t *testing.T, path string, lines ...string) {
	t.Helper()
	c.Write(t, path, strings.Join(lines, "\n"))
}

func (c *TestContext) Cwd(t *testing.T) string {
	t.Helper()
	lines, err := c.shell.Run("pwd")
	require.NoError(t, err)
	require.Len(t, lines, 1, "unexpected output for 'pwd'")
	return lines[0]
}

func (c *TestContext) Cat(t *testing.T, path string) string {
	t.Helper()
	lines, err := c.shell.Run("cat " + strconv.Quote(path))
	require.NoError(t, err)
	return strings.Join(lines, "\n")
}

func (c *TestContext) Ls(t *testing.T, path string) []string {
	t.Helper()
	lines, err := c.shell.Run("ls -1 " + strconv.Quote(path))
	require.NoError(t, err)
	return lines
}

func (c *TestContext) AssertExist(t *testing.T, path string) {
	t.Helper()
	quotedPath := strconv.Quote(path)
	_, err := c.shell.Run("test -e " + strconv.Quote(quotedPath))
	require.NoError(t, err, "expected file %s to exist", quotedPath)
}

func (c *TestContext) AssertContains(t *testing.T, path, expected string) {
	t.Helper()
	value := c.Cat(t, path)
	require.Equal(t, expected, value, "expected file %s to contain %s", strconv.Quote(path), strconv.Quote(expected))
}

func (c *TestContext) GetEnv(t *testing.T, name string) string {
	t.Helper()
	lines, err := c.shell.Run("echo ${" + name + "}")
	require.NoError(t, err)
	require.Len(t, lines, 1)
	return lines[0]
}

func (c *TestContext) Cd(t *testing.T, path string) []string {
	t.Helper()
	lines, err := c.shell.Run("cd " + strconv.Quote(path))
	require.NoError(t, err)
	return lines
}

func (c *TestContext) debugLine(format string, a ...any) {
	if c.debug {
		format = strings.TrimSuffix(format, "\n") + "\n"
		fmt.Printf(format, a...)
	}
}
