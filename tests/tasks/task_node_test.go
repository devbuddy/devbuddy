package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

func Test_Task_Node(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	harness.NewProject(t, c,
		`up:`,
		`- node: '10.15.0'`,
	)

	lines := c.Run(t, "bud up", harness.Timeout(2*time.Minute))
	harness.OutputContains(t, lines, "NodeJS (10.15.0)")
	harness.OutputContains(t, lines, "activated: node[10.15.0]")

	lines = c.Run(t, "node -v")
	harness.OutputEqual(t, lines, "v10.15.0")
}

func Test_Task_Node_Npm_Install(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	harness.NewProject(t, c,
		`up:`,
		`- node: '10.15.0'`,
	)

	c.Write(t, "package.json", `{"dependencies": {"testpackage": "1.0.0"}}`)

	c.Run(t, "bud up", harness.Timeout(2*time.Minute))

	lines := c.Run(t, "npm list")
	harness.OutputContains(t, lines, "testpackage@1.0.0")
}
