package integration

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/tests/context"
)

func CreateContext(t *testing.T) *context.TestContext {
	t.Helper()

	c, err := context.New(config)
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

func OutputContains(t *testing.T, lines []string, subStrings ...string) {
	t.Helper()

	text := strings.Join(lines, "\n")
	text = context.StripAnsi(text)

	for _, subString := range subStrings {
		require.Contains(t, text, subString)
	}
}

func OutputNotContain(t *testing.T, lines []string, subStrings ...string) {
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
	name := fmt.Sprintf("project-%x", rand.Int31())

	p := Project{
		c:    c,
		Path: "/home/tester/src/github.com/orgname/" + name,
	}

	c.Run(t, "mkdir -p "+p.Path)

	p.WriteDevYml(t, devYmlLines...)

	return p
}

func (p *Project) WriteDevYml(t *testing.T, devYmlLines ...string) {
	path := p.Path + "/dev.yml"
	p.c.Write(t, path, strings.Join(devYmlLines, "\n"))
}
