package harness

import (
	stdcontext "context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	testcontext "github.com/devbuddy/devbuddy/tests/context"
	"github.com/devbuddy/devbuddy/tests/internal/shellharness"
)

type PipeContext struct {
	runner                 *shellharness.Runner
	workspaceHostPath      string
	workspaceContainerPath string
}

type PipeProject struct {
	c    *PipeContext
	Path string
}

func CreatePipeContext(t *testing.T) *PipeContext {
	t.Helper()

	testConfig := config
	workspaceHostPath, err := os.MkdirTemp("", "devbuddy-test-workspace-*")
	require.NoError(t, err)
	err = os.Chmod(workspaceHostPath, 0777)
	require.NoError(t, err)

	workspaceContainerPath := "/home/tester/src/github.com"
	shellPath, shellArgs := shellCommand(t, testConfig.ShellName)
	dockerCommand := dockerShellCommand(testConfig, shellPath, shellArgs, workspaceHostPath, workspaceContainerPath)

	// Start the same Docker test image as the PTY harness, but without "-t".
	// Commands are driven through shellharness' sentinel protocol instead of
	// waiting for prompts from an interactive terminal.
	runner, err := shellharness.Start(dockerCommand[0], dockerCommand[1:]...)
	require.NoError(t, err)

	c := &PipeContext{
		runner:                 runner,
		workspaceHostPath:      workspaceHostPath,
		workspaceContainerPath: workspaceContainerPath,
	}
	c.Run(t, "umask 000")
	OutputEqual(t, c.Run(t, "echo $IN_DOCKER"), "yes")

	t.Cleanup(func() {
		err := c.Close()
		require.NoError(t, err)
		_ = os.RemoveAll(workspaceHostPath)
	})

	return c
}

func CreatePipeContextAndInit(t *testing.T) *PipeContext {
	t.Helper()

	c := CreatePipeContext(t)
	OutputEqual(t, c.Run(t, `eval "$(bud --shell-init)"`))
	return c
}

func CreatePipeProject(t *testing.T, c *PipeContext, devYmlLines ...string) PipeProject {
	t.Helper()

	name := fmt.Sprintf("project-%x", rand.Int32())
	p := PipeProject{
		c:    c,
		Path: "/home/tester/src/github.com/orgname/" + name,
	}
	p.WriteDevYml(t, devYmlLines...)
	return p
}

func (p *PipeProject) WriteDevYml(t *testing.T, devYmlLines ...string) {
	t.Helper()

	p.c.Write(t, p.Path+"/dev.yml", strings.Join(devYmlLines, "\n"))
}

func (c *PipeContext) Close() error {
	return c.runner.Close()
}

func (c *PipeContext) Run(t *testing.T, cmd string) []string {
	t.Helper()

	ctx, cancel := stdcontext.WithTimeout(stdcontext.Background(), 10*time.Second)
	defer cancel()

	result, err := c.runner.Run(ctx, cmd)
	require.NoError(t, err, "running command: %q", cmd)
	require.Equal(t, 0, result.ExitCode, "running command: %q. output:\n%s", cmd, strings.Join(result.Lines, "\n"))
	return testcontext.StripAnsiSlice(result.Lines)
}

func (c *PipeContext) Write(t *testing.T, containerPath, content string) {
	t.Helper()

	hostPath := c.hostPath(t, containerPath)
	err := os.MkdirAll(filepath.Dir(hostPath), 0755)
	require.NoError(t, err)
	err = os.Chmod(filepath.Dir(hostPath), 0777)
	require.NoError(t, err)
	err = os.WriteFile(hostPath, []byte(content), 0644)
	require.NoError(t, err)
}

func (c *PipeContext) Cd(t *testing.T, path string) []string {
	t.Helper()
	return c.Run(t, "cd "+strconv.Quote(path))
}

func (c *PipeContext) GetEnv(t *testing.T, name string) string {
	t.Helper()
	lines := c.Run(t, "echo ${"+name+"}")
	require.NotEmpty(t, lines)
	return lines[len(lines)-1]
}

func (c *PipeContext) hostPath(t *testing.T, containerPath string) string {
	t.Helper()

	absolutePath := path.Clean(containerPath)
	workspaceRoot := path.Clean(c.workspaceContainerPath)
	require.True(t, absolutePath == workspaceRoot || strings.HasPrefix(absolutePath, workspaceRoot+"/"), "path %q is outside workspace %q", containerPath, workspaceRoot)

	relPath := "."
	if absolutePath != workspaceRoot {
		relPath = strings.TrimPrefix(absolutePath, workspaceRoot+"/")
	}
	return filepath.Join(c.workspaceHostPath, filepath.FromSlash(relPath))
}

func shellCommand(t *testing.T, shellName string) (string, []string) {
	t.Helper()

	switch shellName {
	case "bash":
		return "/bin/bash", []string{"--noprofile", "--norc"}
	case "zsh":
		return "/bin/zsh", []string{"--no-globalrcs", "--no-rcs", "--no-zle", "--no-promptcr"}
	default:
		require.Failf(t, "unknown shell", "unknown shell: %s", shellName)
		return "", nil
	}
}

func dockerShellCommand(config testcontext.Config, shellPath string, shellArgs []string, workspaceHostPath, workspaceContainerPath string) []string {
	dockerExec := "docker"
	cmd := exec.Command("docker", "-v")
	if cmd.Run() != nil {
		dockerExec = "podman"
	}

	dockerCommand := []string{
		dockerExec, "run",
		"-i",
		"-v", config.BinaryPath + ":/usr/local/bin/bud",
		"-v", workspaceHostPath + ":" + workspaceContainerPath,
		"-e", "IN_DOCKER=yes",
		"--rm",
		"--entrypoint", shellPath,
		config.DockerImage,
	}
	return append(dockerCommand, shellArgs...)
}
