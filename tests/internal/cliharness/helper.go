package cliharness

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	testcontext "github.com/devbuddy/devbuddy/tests/context"
)

type Context struct {
	workDir string
	homeDir string
	env     []string
}

type RunOptions struct {
	ExitCode int
	Timeout  time.Duration
}

type RunOption func(*RunOptions)

func ExitCode(code int) RunOption {
	return func(o *RunOptions) {
		o.ExitCode = code
	}
}

func Timeout(timeout time.Duration) RunOption {
	return func(o *RunOptions) {
		o.Timeout = timeout
	}
}

func NewContext(t *testing.T) *Context {
	t.Helper()

	workDir := canonicalTempDir(t)
	homeDir := canonicalTempDir(t)

	pathValue := filepath.Dir(binaryPath) + string(os.PathListSeparator) + os.Getenv("PATH")
	env := append(os.Environ(),
		"HOME="+homeDir,
		"PATH="+pathValue,
	)

	return &Context{
		workDir: workDir,
		homeDir: homeDir,
		env:     env,
	}
}

func (c *Context) Run(t *testing.T, command string, optFns ...RunOption) []string {
	t.Helper()

	opt := RunOptions{
		ExitCode: 0,
		Timeout:  10 * time.Second,
	}
	for _, fn := range optFns {
		fn(&opt)
	}

	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
	defer cancel()

	finalizerFile := filepath.Join(t.TempDir(), "bud-finalizer")
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = c.workDir
	cmd.Env = append(append([]string{}, c.env...), "BUD_FINALIZER_FILE="+finalizerFile)

	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
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

	return lines(output)
}

func (c *Context) Write(t *testing.T, path, content string) {
	t.Helper()
	hostPath := c.resolvePath(path)
	err := os.MkdirAll(filepath.Dir(hostPath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(hostPath, []byte(content), 0644)
	require.NoError(t, err)
}

func (c *Context) WriteLines(t *testing.T, path string, lines ...string) {
	t.Helper()
	c.Write(t, path, strings.Join(lines, "\n"))
}

func (c *Context) Cwd(t *testing.T) string {
	t.Helper()
	return c.workDir
}

func (c *Context) Cd(t *testing.T, path string) []string {
	t.Helper()
	c.workDir = c.resolvePath(path)
	err := os.MkdirAll(c.workDir, 0755)
	require.NoError(t, err)
	return nil
}

func (c *Context) Setenv(name, value string) {
	c.setenv(name, value)
}

func (c *Context) PrependPath(path string) {
	c.setenv("PATH", path+string(os.PathListSeparator)+c.getenv("PATH"))
}

func (c *Context) Cat(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(c.resolvePath(path))
	require.NoError(t, err)
	return strings.TrimSuffix(string(content), "\n")
}

func (c *Context) Ls(t *testing.T, path string) []string {
	t.Helper()
	entries, err := os.ReadDir(c.resolvePath(path))
	require.NoError(t, err)

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

func (c *Context) AssertContains(t *testing.T, path, expected string) {
	t.Helper()
	require.Equal(t, expected, c.Cat(t, path), "expected file %q to contain %q", path, expected)
}

func (c *Context) Path(path string) string {
	return c.resolvePath(path)
}

func (c *Context) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(c.workDir, path)
}

func (c *Context) getenv(name string) string {
	prefix := name + "="
	for i := len(c.env) - 1; i >= 0; i-- {
		if strings.HasPrefix(c.env[i], prefix) {
			return strings.TrimPrefix(c.env[i], prefix)
		}
	}
	return ""
}

func (c *Context) setenv(name, value string) {
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

type Project struct {
	c    *Context
	Path string
}

func CreateContext(t *testing.T) *Context {
	t.Helper()
	return NewContext(t)
}

func CreateContextAndInit(t *testing.T) *Context {
	t.Helper()
	return CreateContext(t)
}

func CreateContextAndProject(t *testing.T, devYmlLines ...string) (*Context, Project) {
	t.Helper()

	c := CreateContext(t)
	p := CreateProject(t, c, devYmlLines...)
	c.Cd(t, p.Path)
	return c, p
}

func CreateProject(t *testing.T, c *Context, devYmlLines ...string) Project {
	t.Helper()

	name := fmt.Sprintf("project-%x", rand.Int32())
	p := Project{
		c:    c,
		Path: filepath.Join(c.homeDir, "src", "github.com", "orgname", name),
	}

	p.WriteDevYml(t, devYmlLines...)
	return p
}

func (p *Project) WriteDevYml(t *testing.T, devYmlLines ...string) {
	t.Helper()
	p.c.Write(t, filepath.Join(p.Path, "dev.yml"), strings.Join(devYmlLines, "\n"))
}

func OutputContains(t *testing.T, output []string, subStrings ...string) {
	t.Helper()

	text := testcontext.StripAnsi(strings.Join(output, "\n"))
	for _, subString := range subStrings {
		require.Contains(t, text, subString)
	}
}

func OutputNotContains(t *testing.T, output []string, subStrings ...string) {
	t.Helper()

	text := testcontext.StripAnsi(strings.Join(output, "\n"))
	for _, subString := range subStrings {
		require.NotContains(t, text, subString)
	}
}

func OutputEqual(t *testing.T, output []string, expectedLines ...string) {
	t.Helper()
	require.Equal(t, expectedLines, output)
}

func lines(output []byte) []string {
	text := strings.ReplaceAll(string(output), "\r", "")
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(testcontext.StripAnsi(text), "\n")
}

func canonicalTempDir(t *testing.T) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	return dir
}
