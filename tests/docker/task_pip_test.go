package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Task_Pip(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- pip : [requirements.txt]`,
	)
	c.Cd(t, p.Path)

	c.Write(t, "requirements.txt", "pkginfo==1.9.6\n")

	lines := c.Run(t, "bud up", context.Timeout(2*time.Minute))
	OutputContains(t, lines, "activated: python 3.9.0")

	lines = c.Run(t, "pip freeze")
	OutputContains(t, lines, "pkginfo==1.9.6")
}
