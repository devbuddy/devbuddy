package harness

import (
	stdcontext "context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/tests/internal/context"
)

var cliBinaryPath string

// SetupCLI builds the host `bud` binary and prepares the CLI harness. It is
// intended to be called from a test package's TestMain before m.Run().
func SetupCLI() error {
	path, err := BuildHostBinary()
	if err != nil {
		return fmt.Errorf("building host bud binary: %w", err)
	}
	cliBinaryPath = path
	return nil
}

// CLIContext runs `bud` as a host subprocess in a temporary HOME directory.
type CLIContext struct {
	workDir string
	homeDir string
	env     []string
}

// RunOption configures a single CLIContext.Run invocation.
type RunOption func(*runOptions)

type runOptions struct {
	ExitCode int
	Timeout  time.Duration
}

// ExitCode sets the expected exit code for a CLIContext.Run call.
func ExitCode(code int) RunOption {
	return func(o *runOptions) { o.ExitCode = code }
}

// Timeout sets the per-command timeout for a CLIContext.Run call.
func Timeout(timeout time.Duration) RunOption {
	return func(o *runOptions) { o.Timeout = timeout }
}

// NewCLI creates a host CLI context with a fresh temp HOME and work dir.
func NewCLI(t *testing.T) *CLIContext {
	t.Helper()
	workDir := canonicalTempDir(t)
	homeDir := canonicalTempDir(t)
	pathValue := filepath.Dir(cliBinaryPath) + string(os.PathListSeparator) + os.Getenv("PATH")
	env := append(os.Environ(),
		"HOME="+homeDir,
		"PATH="+pathValue,
	)
	return &CLIContext{workDir: workDir, homeDir: homeDir, env: env}
}

// NewCLIProject creates a project directory with the given dev.yml content,
// cd's into it, and returns its absolute path.
func NewCLIProject(t *testing.T, c *CLIContext, devYmlLines ...string) string {
	t.Helper()
	name := fmt.Sprintf("project-%x", rand.Int32())
	projectPath := filepath.Join(c.homeDir, "src", "github.com", "orgname", name)
	c.WriteLines(t, filepath.Join(projectPath, "dev.yml"), devYmlLines...)
	c.Cd(t, projectPath)
	return projectPath
}

// Run executes a shell command in the context's working directory.
func (c *CLIContext) Run(t *testing.T, command string, optFns ...RunOption) []string {
	t.Helper()

	opt := runOptions{ExitCode: 0, Timeout: 10 * time.Second}
	for _, fn := range optFns {
		fn(&opt)
	}

	ctx, cancel := stdcontext.WithTimeout(stdcontext.Background(), opt.Timeout)
	defer cancel()

	finalizerFile := filepath.Join(t.TempDir(), "bud-finalizer")
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = c.workDir
	cmd.Env = append(append([]string{}, c.env...), "BUD_FINALIZER_FILE="+finalizerFile)

	output, err := cmd.CombinedOutput()
	if ctx.Err() == stdcontext.DeadlineExceeded {
		require.Failf(t, "command timed out", "timed out after %s running %q. output:\n%s", opt.Timeout, command, string(output))
	}

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			require.NoError(t, err, "running command: %q", command)
		}
	}
	require.Equal(t, opt.ExitCode, exitCode, "running command %q. output:\n%s", command, string(output))

	return splitLines(output)
}

// Write writes a file at path (resolved against the work dir if relative).
func (c *CLIContext) Write(t *testing.T, path, content string) {
	t.Helper()
	hostPath := c.resolvePath(path)
	require.NoError(t, os.MkdirAll(filepath.Dir(hostPath), 0755))
	require.NoError(t, os.WriteFile(hostPath, []byte(content), 0644))
}

// WriteLines writes content joined by newlines.
func (c *CLIContext) WriteLines(t *testing.T, path string, lines ...string) {
	t.Helper()
	c.Write(t, path, strings.Join(lines, "\n"))
}

// Cwd returns the current working directory.
func (c *CLIContext) Cwd(t *testing.T) string {
	t.Helper()
	return c.workDir
}

// Cd changes the working directory (creating it if needed).
func (c *CLIContext) Cd(t *testing.T, path string) []string {
	t.Helper()
	c.workDir = c.resolvePath(path)
	require.NoError(t, os.MkdirAll(c.workDir, 0755))
	return nil
}

// Setenv sets an env var for subsequent Run calls.
func (c *CLIContext) Setenv(name, value string) {
	c.setenv(name, value)
}

// PrependPath prepends a directory to PATH.
func (c *CLIContext) PrependPath(path string) {
	c.setenv("PATH", path+string(os.PathListSeparator)+c.getenv("PATH"))
}

// Cat reads a file from the context, trimming a trailing newline.
func (c *CLIContext) Cat(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(c.resolvePath(path))
	require.NoError(t, err)
	return strings.TrimSuffix(string(b), "\n")
}

// Ls lists names of entries under path.
func (c *CLIContext) Ls(t *testing.T, path string) []string {
	t.Helper()
	entries, err := os.ReadDir(c.resolvePath(path))
	require.NoError(t, err)
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

// AssertContains asserts that the file equals the expected string.
func (c *CLIContext) AssertContains(t *testing.T, path, expected string) {
	t.Helper()
	require.Equal(t, expected, c.Cat(t, path), "expected file %q to contain %q", path, expected)
}

// Path resolves a path against the context's work dir if relative.
func (c *CLIContext) Path(path string) string {
	return c.resolvePath(path)
}

func (c *CLIContext) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(c.workDir, path)
}

func (c *CLIContext) getenv(name string) string {
	prefix := name + "="
	for i := len(c.env) - 1; i >= 0; i-- {
		if v, ok := strings.CutPrefix(c.env[i], prefix); ok {
			return v
		}
	}
	return ""
}

func (c *CLIContext) setenv(name, value string) {
	prefix := name + "="
	filtered := c.env[:0]
	for _, entry := range c.env {
		if !strings.HasPrefix(entry, prefix) {
			filtered = append(filtered, entry)
		}
	}
	c.env = filtered
	c.env = append(c.env, prefix+value)
}

func splitLines(output []byte) []string {
	text := strings.ReplaceAll(string(output), "\r", "")
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(context.StripAnsi(text), "\n")
}

func canonicalTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	return dir
}
