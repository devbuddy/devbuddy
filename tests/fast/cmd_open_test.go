package integration

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func installFakeXdgOpen(t *testing.T, c *Context) string {
	t.Helper()

	binDir := c.Path("bin")
	outputPath := c.Path("devbuddy-open-url")
	script := `#!/bin/sh
printf "%s\n" "$1" > ` + outputPath + `
`
	c.Write(t, binDir+"/xdg-open", script)
	c.Write(t, binDir+"/open", script)
	c.Run(t, "chmod +x "+binDir+"/xdg-open")
	c.Run(t, "chmod +x "+binDir+"/open")
	c.PrependPath(binDir)

	return outputPath
}

func waitAndReadOpenedURL(t *testing.T, c *Context, path string) string {
	t.Helper()

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		_, err := os.ReadFile(path)
		if err == nil {
			return c.Run(t, "cat "+path)[0]
		}
		time.Sleep(10 * time.Millisecond)
	}

	require.Failf(t, "opened URL was not recorded", "expected %s to be written", path)
	return ""
}

func Test_Cmd_Open_CustomLink_FuzzyMatch(t *testing.T) {
	c, p := CreateContextAndProject(t,
		`open:`,
		`  staging: https://staging.example.com`,
		`  docs: https://docs.example.com`,
	)
	outputPath := installFakeXdgOpen(t, c)
	c.Cd(t, p.Path)

	c.Run(t, "bud open stg")

	openedURL := waitAndReadOpenedURL(t, c, outputPath)
	OutputEqual(t, []string{openedURL}, "https://staging.example.com")
}

func Test_Cmd_Open_NoArgOpensGithub(t *testing.T) {
	c, p := CreateContextAndProject(t,
		`open:`,
		`  docs: https://docs.example.com`,
	)
	c.Cd(t, p.Path)

	// No git remote configured, so it fails
	lines := c.Run(t, "bud open", ExitCode(1))
	OutputContains(t, lines, "failed to get the origin remote url")
}
