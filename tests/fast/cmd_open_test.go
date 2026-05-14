package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
	"github.com/stretchr/testify/require"
)

func installFakeXdgOpen(t *testing.T, c *context.TestContext) {
	t.Helper()

	c.Run(t, "mkdir -p /tmp/devbuddy-test-bin")
	c.Write(t, "/tmp/devbuddy-test-bin/xdg-open", `#!/bin/sh
printf "%s\n" "$1" > /tmp/devbuddy-open-url
`)
	c.Run(t, "chmod +x /tmp/devbuddy-test-bin/xdg-open")
	c.Run(t, `export PATH="/tmp/devbuddy-test-bin:${PATH}"`)
}

func waitAndReadOpenedURL(t *testing.T, c *context.TestContext) string {
	t.Helper()

	lines := c.Run(t, "cat /tmp/devbuddy-open-url", context.Timeout(15*time.Second))
	require.Len(t, lines, 1)
	return lines[0]
}

func Test_Cmd_Open_CustomLink_FuzzyMatch(t *testing.T) {
	c, p := CreateContextAndProject(t,
		`open:`,
		`  staging: https://staging.example.com`,
		`  docs: https://docs.example.com`,
	)
	installFakeXdgOpen(t, c)
	c.Cd(t, p.Path)

	c.Run(t, "bud open stg")

	openedURL := waitAndReadOpenedURL(t, c)
	OutputEqual(t, []string{openedURL}, "https://staging.example.com")
}

func Test_Cmd_Open_DefaultWhenSingleLink(t *testing.T) {
	c, p := CreateContextAndProject(t,
		`open:`,
		`  docs: https://docs.example.com`,
	)
	installFakeXdgOpen(t, c)
	c.Cd(t, p.Path)

	c.Run(t, "bud open")

	openedURL := waitAndReadOpenedURL(t, c)
	OutputEqual(t, []string{openedURL}, "https://docs.example.com")
}

func Test_Cmd_Open_RequiresNameWhenMultipleLinks(t *testing.T) {
	c, p := CreateContextAndProject(t,
		`open:`,
		`  staging: https://staging.example.com`,
		`  docs: https://docs.example.com`,
	)
	c.Cd(t, p.Path)

	lines := c.Run(t, "bud open", context.ExitCode(1))
	OutputEqual(t, lines, "Error: which link should I open?")
}
