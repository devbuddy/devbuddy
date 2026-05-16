package harness

import (
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/tests/internal/context"
)

var dockerConfig context.Config

// SetupDocker loads the test config and builds the Linux `bud` binary. It is
// intended to be called from a test package's TestMain before m.Run().
func SetupDocker() error {
	cfg, err := context.LoadConfig()
	if err != nil {
		return fmt.Errorf("loading test config: %w", err)
	}

	binaryPath, err := BuildLinuxBinary()
	if err != nil {
		return fmt.Errorf("building linux bud binary: %w", err)
	}
	cfg.BinaryPath = binaryPath

	dockerConfig = cfg
	return nil
}

// NewDocker creates a non-PTY Docker context. Use this when the test does not
// need a controlling terminal.
func NewDocker(t *testing.T) *context.TestContext {
	t.Helper()
	return newDockerContext(t, false)
}

// NewDockerPTY creates a PTY-backed Docker context.
func NewDockerPTY(t *testing.T) *context.TestContext {
	t.Helper()
	return newDockerContext(t, true)
}

// NewDockerInit creates a non-PTY Docker context and evaluates `bud --shell-init`.
func NewDockerInit(t *testing.T) *context.TestContext {
	t.Helper()
	c := NewDocker(t)
	initShell(t, c)
	return c
}

// NewDockerPTYInit creates a PTY Docker context and evaluates `bud --shell-init`.
func NewDockerPTYInit(t *testing.T) *context.TestContext {
	t.Helper()
	c := NewDockerPTY(t)
	initShell(t, c)
	return c
}

// NewDockerProject creates a project directory with the given dev.yml content,
// cd's into it, and returns its absolute path.
func NewDockerProject(t *testing.T, c *context.TestContext, devYmlLines ...string) string {
	t.Helper()
	name := fmt.Sprintf("project-%x", rand.Int32())
	projectPath := "/home/tester/src/github.com/orgname/" + name
	c.WriteLines(t, projectPath+"/dev.yml", devYmlLines...)
	c.Cd(t, projectPath)
	return projectPath
}

func newDockerContext(t *testing.T, usePTY bool) *context.TestContext {
	t.Helper()

	testConfig := dockerConfig
	workspaceHostPath, err := os.MkdirTemp("", "devbuddy-test-workspace-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(workspaceHostPath)
	})
	require.NoError(t, os.Chmod(workspaceHostPath, 0777))

	testConfig.WorkspaceHostPath = workspaceHostPath
	testConfig.WorkspaceContainerPath = "/home/tester/src/github.com"
	testConfig.UsePTY = usePTY

	c, err := context.New(testConfig)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, c.Close())
	})
	return c
}

func initShell(t *testing.T, c *context.TestContext) {
	t.Helper()
	output := c.Run(t, `eval "$(bud --shell-init)"`)
	require.Len(t, output, 0)
}
