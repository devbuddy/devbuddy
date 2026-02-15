package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Task_Node(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up:`,
		`- node: '10.15.0'`,
	)
	c.Cd(t, p.Path)

	lines := c.Run(t, "bud up", context.Timeout(2*time.Minute))
	OutputContains(t, lines, "NodeJS (10.15.0)")
	OutputContains(t, lines, "activated: node 10.15.0")

	lines = c.Run(t, "node -v")
	OutputEqual(t, lines, "v10.15.0")
}

func Test_Task_Node_Npm_Install(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up:`,
		`- node: '10.15.0'`,
	)
	c.Cd(t, p.Path)

	c.Write(t, "package.json", `{"dependencies": {"testpackage": "1.0.0"}}`)

	c.Run(t, "bud up", context.Timeout(2*time.Minute))

	lines := c.Run(t, "npm list")
	OutputContains(t, lines, "testpackage@1.0.0")
}
