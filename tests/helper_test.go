package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/tests/context"
)

var config context.Config // Initialized by TestMain()

func CreateContext(t *testing.T) *context.TestContext {
	t.Helper()

	c, err := context.New(t, config)
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
	output := c.Run(`eval "$(bud --shell-init)"`)
	require.Len(t, output, 0)
	return c
}

func OutputContains(t *testing.T, lines []string, subStrings ...string) {
	t.Helper()

	text := strings.Join(lines, "\n")

	for _, subString := range subStrings {
		if !strings.Contains(text, subString) {
			t.Fatalf("Substring %q was not found in:\n%s", subString, text)
		}
	}
}

func OutputNotContain(t *testing.T, lines []string, subStrings ...string) {
	t.Helper()

	text := strings.Join(lines, "\n")

	for _, subString := range subStrings {
		if strings.Contains(text, subString) {
			t.Fatalf("Substring %q was found in:\n%s", subString, text)
		}
	}
}

func OutputEqual(t *testing.T, lines []string, expectedLines ...string) {
	t.Helper()
	require.Equal(t, expectedLines, lines)
}

type Project struct {
	Name string
	Path string
}

func CreateProject(c *context.TestContext, name string, devYmlLines ...string) Project {
	projectPath := "/home/tester/src/github.com/orgname/" + name
	c.Run("mkdir -p " + projectPath)

	path := projectPath + "/dev.yml"
	c.Write(path, strings.Join(devYmlLines, "\n"))
	c.Run("bud cd " + name)

	return Project{name, projectPath}
}

func (p *Project) UpdateDevYml(c *context.TestContext, devYmlLines ...string) {
	path := p.Path + "/dev.yml"
	c.Write(path, strings.Join(devYmlLines, "\n"))
}
