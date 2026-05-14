package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Task_Pipfile(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- pipfile`,
	)
	c.Cd(t, p.Path)

	c.Write(t, "Pipfile", "[packages]\n\"pkginfo\" = \"==1.9.6\"\n")

	c.Run(t, "bud up", context.Timeout(2*time.Minute))

	lines := c.Run(t, "pip freeze")
	OutputContains(t, lines, "pkginfo==1.9.6")
}
