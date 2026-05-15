package harness

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/tests/context"
)

func CreateContext(t *testing.T) *context.TestContext {
	t.Helper()

	testConfig := config
	var err error
	testConfig.WorkspaceHostPath, err = os.MkdirTemp("", "devbuddy-test-workspace-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(testConfig.WorkspaceHostPath)
	})
	err = os.Chmod(testConfig.WorkspaceHostPath, 0777)
	require.NoError(t, err)
	testConfig.WorkspaceContainerPath = "/home/tester/src/github.com"

	c, err := context.New(testConfig)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := c.Close()
		require.NoError(t, err)
	})

	return c
}

func CreateContextAndInit(t *testing.T) *context.TestContext {
	t.Helper()

	c := CreateContext(t)
	output := c.Run(t, `eval "$(bud --shell-init)"`)
	require.Len(t, output, 0)
	return c
}

// CreateContextAndProject creates an initialized context, a project with the
// given dev.yml content, and cd's into the project directory. This combines the
// common 3-line boilerplate of CreateContextAndInit + CreateProject + Cd.
func CreateContextAndProject(t *testing.T, devYmlLines ...string) (*context.TestContext, Project) {
	t.Helper()

	c := CreateContextAndInit(t)
	p := CreateProject(t, c, devYmlLines...)
	c.Cd(t, p.Path)
	return c, p
}

func OutputContains(t *testing.T, lines []string, subStrings ...string) {
	t.Helper()

	text := strings.Join(lines, "\n")
	text = context.StripAnsi(text)

	for _, subString := range subStrings {
		require.Contains(t, text, subString)
	}
}

func OutputNotContains(t *testing.T, lines []string, subStrings ...string) {
	t.Helper()

	text := strings.Join(lines, "\n")
	text = context.StripAnsi(text)

	for _, subString := range subStrings {
		require.NotContains(t, text, subString)
	}
}

func OutputEqual(t *testing.T, lines []string, expectedLines ...string) {
	t.Helper()
	require.Equal(t, expectedLines, lines)
}

type Project struct {
	c    *context.TestContext
	Path string
}

func CreateProject(t *testing.T, c *context.TestContext, devYmlLines ...string) Project {
	t.Helper()

	name := fmt.Sprintf("project-%x", rand.Int32())

	p := Project{
		c:    c,
		Path: "/home/tester/src/github.com/orgname/" + name,
	}

	p.WriteDevYml(t, devYmlLines...)

	return p
}

func (p *Project) WriteDevYml(t *testing.T, devYmlLines ...string) {
	t.Helper()

	path := p.Path + "/dev.yml"
	p.c.Write(t, path, strings.Join(devYmlLines, "\n"))
}
