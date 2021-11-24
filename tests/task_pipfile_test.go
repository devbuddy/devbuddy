package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Task_Pipfile(t *testing.T) {
	c := CreateContextAndInit(t)

	CreateProject(c, "project",
		`up:`,
		`- python: 3.9.0`,
		`- pipfile`,
	)

	c.Write("Pipfile", `[packages]\n"pyreleaser" = "==0.5.2"\n`)

	c.Run("bud up", context.Timeout(2*time.Minute))

	lines := c.Run("pip freeze")
	OutputContains(t, lines, "pyreleaser==0.5.2")
}
