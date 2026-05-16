package integration

import (
	"os"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
	"github.com/stretchr/testify/require"
)

func installFakeXdgOpen(t *testing.T, c *harness.CLIContext) string {
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

func waitAndReadOpenedURL(t *testing.T, c *harness.CLIContext, path string) string {
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
	c := harness.NewCLI(t)
	harness.NewCLIProject(t, c,
		`open:`,
		`  staging: https://staging.example.com`,
		`  docs: https://docs.example.com`,
	)
	outputPath := installFakeXdgOpen(t, c)

	c.Run(t, "bud open stg")

	openedURL := waitAndReadOpenedURL(t, c, outputPath)
	harness.OutputEqual(t, []string{openedURL}, "https://staging.example.com")
}

func Test_Cmd_Open_NoArgOpensGithub(t *testing.T) {
	c := harness.NewCLI(t)
	harness.NewCLIProject(t, c,
		`open:`,
		`  docs: https://docs.example.com`,
	)

	// No git remote configured, so it fails
	lines := c.Run(t, "bud open", harness.ExitCode(1))
	harness.OutputContains(t, lines, "failed to get the origin remote url")
}
