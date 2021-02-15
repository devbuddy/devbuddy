package context

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/integration/context/expect"
	"github.com/stretchr/testify/require"
)

type TestContext struct {
	expect *expect.ExpectProcess
	shell  interface {
		Run(string) ([]string, error)
	}
	t     *testing.T
	debug bool
}

func New(config Config, t *testing.T) (*TestContext, error) {
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
		panic("unknown shell " + config.ShellName)
	}

	dockerCommand := []string{
		"docker", "run",
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
		t:      t,
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

func (e *TestContext) Close() error {
	return e.expect.Stop()
}

func (e *TestContext) Run(cmd string, optFns ...runOptionsFn) []string {
	lines, err := e.run(cmd, optFns...)
	require.NoError(e.t, err, "running command: %q", cmd)
	return lines
}

func (e *TestContext) run(cmd string, optFns ...runOptionsFn) ([]string, error) {
	opt := buildRunOptions(optFns)

	e.debugLine("Running command %q", cmd)
	e.debugLine("Options: %+v", opt)

	var lines []string
	var err error

	done := make(chan bool)
	go func() {
		lines, err = e.shell.Run(cmd)
		close(done)
	}()

	select {
	case <-done:
		e.debugLine("command completed")
	case <-time.After(opt.timeout):
		return nil, fmt.Errorf("timed out after %s", opt.timeout)
	}

	codeLines, err := e.shell.Run("echo $?")
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
		exert := "no output"
		if len(lines) > 0 {
			exert = lines[0]
		}
		return nil, fmt.Errorf("exit code %d. first output line: %s", exitCode, exert)
	}

	return StripAnsiSlice(lines), nil
}

func (e *TestContext) Write(path, content string) {
	_, err := e.shell.Run(fmt.Sprintf("echo -e %q > %q", content, path))
	require.NoError(e.t, err)
}

func (e *TestContext) Cwd() string {
	lines, err := e.shell.Run("pwd")
	require.NoError(e.t, err)
	require.Len(e.t, lines, 1, "unexpected output for 'pwd'")
	return lines[0]
}

func (e *TestContext) Cat(path string) string {
	lines, err := e.shell.Run("cat " + strconv.Quote(path))
	require.NoError(e.t, err)
	return strings.Join(lines, "\n")
}

func (e *TestContext) Ls(path string) []string {
	lines, err := e.shell.Run("ls -1 " + strconv.Quote(path))
	require.NoError(e.t, err)
	return lines
}

func (e *TestContext) AssertExist(path string) {
	quotedPath := strconv.Quote(path)
	_, err := e.shell.Run("test -e " + strconv.Quote(quotedPath))
	require.NoError(e.t, err, "expected file %s to exist", quotedPath)
}

func (e *TestContext) GetEnv(name string) string {
	lines, err := e.shell.Run("echo ${" + name + "}")
	require.NoError(e.t, err)
	require.Len(e.t, lines, 1)
	return lines[0]
}

func (e *TestContext) Cd(path string) {
	_, err := e.shell.Run("cd " + strconv.Quote(path))
	require.NoError(e.t, err)
}

func (e *TestContext) debugLine(format string, a ...interface{}) {
	if e.debug {
		format = strings.TrimSuffix(format, "\n") + "\n"
		fmt.Printf(format, a...)
	}
}
