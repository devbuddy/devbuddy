package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Cmd_Inspect(t *testing.T) {
	c := CreateContextAndInit(t)

	project := CreateProject(t, c,
		`up:`,
		`- node: '10.15.0'`,
		`- pip: [requirements.txt]`,
	)
	c.Cd(t, project.Path)

	lines := c.Run(t, "bud inspect")
	OutputEqual(t, lines[0:3],
		"Found project at "+project.Path,
		"- Task NodeJS (10.15.0) feature=node:10.15.0 actions=2",
		"- Task Pip (requirements.txt) required_task=python actions=1",
	)
}

func Test_Cmd_Inspect_No_Manifest(t *testing.T) {
	c := CreateContextAndInit(t)

	lines := c.Run(t, "bud inspect", context.ExitCode(1))
	OutputEqual(t, lines, "Error: project not found")
}
