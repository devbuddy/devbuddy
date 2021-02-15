package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Task_Pip(t *testing.T) {
	c := CreateContextAndInit(t)

	CreateProject(t, c, "project",
		`up:`,
		`- python: 3.9.0`,
		`- pip : [requirements.txt]`,
	)

	c.Write("requirements.txt", "pyreleaser==0.5.2\n")

	lines := c.Run("bud up", context.Timeout(2*time.Minute))
	OutputContains(t, lines, "python activated. (3.9.0)")

	lines = c.Run("pip freeze")
	OutputContains(t, lines, "pyreleaser==0.5.2")
}
