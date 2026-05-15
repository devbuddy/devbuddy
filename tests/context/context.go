package context

import (
	stdcontext "context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/tests/internal/expect"
	"github.com/devbuddy/devbuddy/tests/internal/shellharness"
)

type TestContext struct {
	expect *expect.ExpectProcess
	shell  interface {
		Run(string) ([]string, error)
	}
	close                  func() error
	workspaceHostPath      string
	workspaceContainerPath string
	debug                  bool
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
		"-v", config.BinaryPath + ":/usr/local/bin/bud",
		"-e", "PROMPT=##\n",
		"-e", "PS1=##\n",
		"-e", "IN_DOCKER=yes",
		"--rm",
		"-i", // Keep STDIN open even if not attached
	}
	if config.UsePTY {
		dockerCommand = append(dockerCommand, "-t") // Allocate a pseudo-TTY
	}
	if config.WorkspaceHostPath != "" && config.WorkspaceContainerPath != "" {
		dockerCommand = append(dockerCommand, "-v", config.WorkspaceHostPath+":"+config.WorkspaceContainerPath)
	}
	dockerCommand = append(dockerCommand, "--entrypoint", shellPath, config.DockerImage)
	dockerCommand = append(dockerCommand, args...)

	tc := &TestContext{
		workspaceHostPath:      config.WorkspaceHostPath,
		workspaceContainerPath: config.WorkspaceContainerPath,
	}

	if config.UsePTY {
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

		tc.expect = e
		tc.shell = c
		tc.close = e.Stop

		_, err = tc.run("stty -echo")
		if err != nil {
			return nil, fmt.Errorf("disabling echo mode: %w", err)
		}
	} else {
		runner, err := shellharness.Start(dockerCommand[0], dockerCommand[1:]...)
		if err != nil {
			return nil, fmt.Errorf("starting shell runner: %w", err)
		}
		tc.shell = pipeShell{runner: runner}
		tc.close = runner.Close
	}
	tc.debugLine("Shell command: %q", dockerCommand)

	_, err := tc.run("umask 000")
	if err != nil {
		return nil, fmt.Errorf("configuring test shell umask: %w", err)
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
	if c.expect != nil {
		c.expect.Debug = false
	}
}

func (c *TestContext) Debug() {
	c.debug = true
	if c.expect != nil {
		c.expect.Debug = true
	}
}

func (c *TestContext) Close() error {
	c.debugLine("Stopping docker container")
	err := c.close()
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
	var exitCode int

	if shell, ok := c.shell.(interface {
		RunWithExitCode(string, time.Duration) ([]string, int, error)
	}); ok {
		lines, exitCode, err = shell.RunWithExitCode(cmd, opt.timeout)
		if err != nil {
			return nil, err
		}
	} else {
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

		exitCode, err = strconv.Atoi(codeLines[0])
		if err != nil {
			return nil, fmt.Errorf("unexpected exit code %s: %w", codeLines[0], err)
		}
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

	if hostPath, ok := c.hostPath(t, path); ok {
		dirPath := filepath.Dir(hostPath)
		if err := os.MkdirAll(dirPath, 0777); err == nil {
			// chmod each dir from the leaf up to the workspace root so the container
			// user can create siblings. May fail for container-created dirs (different
			// owner on Linux) — best-effort only.
			for current := dirPath; current != c.workspaceHostPath; current = filepath.Dir(current) {
				_ = os.Chmod(current, 0777)
			}
			if err := os.WriteFile(hostPath, []byte(content), 0644); err == nil {
				return
			}
		}
		// Fall through: directory owned by container user (Linux CI) — write via shell.
	}

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
	require.NotEmpty(t, lines, "unexpected output for 'pwd'")
	return lines[len(lines)-1]
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
	_, err := c.shell.Run("test -e " + strconv.Quote(path))
	require.NoError(t, err, "expected file %s to exist", path)
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

func (c *TestContext) hostPath(t *testing.T, containerPath string) (string, bool) {
	t.Helper()

	if c.workspaceHostPath == "" || c.workspaceContainerPath == "" {
		return "", false
	}

	absolutePath := containerPath
	if !path.IsAbs(absolutePath) {
		absolutePath = path.Join(c.Cwd(t), absolutePath)
	}
	absolutePath = path.Clean(absolutePath)

	workspaceRoot := path.Clean(c.workspaceContainerPath)
	if absolutePath != workspaceRoot && !strings.HasPrefix(absolutePath, workspaceRoot+"/") {
		return "", false
	}

	relPath := "."
	if absolutePath != workspaceRoot {
		relPath = strings.TrimPrefix(absolutePath, workspaceRoot+"/")
	}
	return filepath.Join(c.workspaceHostPath, filepath.FromSlash(relPath)), true
}

func (c *TestContext) Send(t *testing.T, text string) {
	t.Helper()
	err := c.expect.Send(text)
	require.NoError(t, err)
}

func (c *TestContext) Expect(t *testing.T, text string) string {
	t.Helper()
	line, err := c.expect.Expect(text)
	require.NoError(t, err)
	return StripAnsi(line)
}

func (c *TestContext) WaitPrompt(t *testing.T) []string {
	t.Helper()
	var output []string
	for {
		line, err := c.expect.Line()
		require.NoError(t, err)

		line = strings.ReplaceAll(line, "\r", "")
		if line == "##\n" {
			return StripAnsiSlice(output)
		}
		output = append(output, strings.TrimSuffix(line, "\n"))
	}
}

func (c *TestContext) debugLine(format string, a ...any) {
	if c.debug {
		format = strings.TrimSuffix(format, "\n") + "\n"
		fmt.Printf(format, a...)
	}
}

type pipeShell struct {
	runner *shellharness.Runner
}

func (s pipeShell) Run(command string) ([]string, error) {
	lines, _, err := s.RunWithExitCode(command, 10*time.Second)
	return lines, err
}

func (s pipeShell) RunWithExitCode(command string, timeout time.Duration) ([]string, int, error) {
	ctx, cancel := stdcontext.WithTimeout(stdcontext.Background(), timeout)
	defer cancel()

	result, err := s.runner.Run(ctx, command)
	if err != nil {
		return nil, 0, err
	}
	return result.Lines, result.ExitCode, nil
}
