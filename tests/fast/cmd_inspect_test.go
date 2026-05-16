package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

func Test_Cmd_Inspect(t *testing.T) {
	c := harness.NewCLI(t)

	projectPath := harness.NewProject(t, c,
		`up:`,
		`- node: '10.15.0'`,
		`- pip: [requirements.txt]`,
	)

	lines := c.Run(t, "bud inspect")
	harness.OutputEqual(t, lines[0:3],
		"Found project at "+projectPath,
		"- Task NodeJS (10.15.0) feature=node:10.15.0 actions=2",
		"- Task Pip (requirements.txt) required_task=python actions=1",
	)
}

func Test_Cmd_Inspect_No_Manifest(t *testing.T) {
	c := harness.NewCLI(t)

	lines := c.Run(t, "bud inspect", harness.ExitCode(1))
	harness.OutputEqual(t, lines, "Error: project not found")
}
