package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

func Test_Task_Pip(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	harness.NewProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- pip : [requirements.txt]`,
	)

	c.Write(t, "requirements.txt", "pkginfo==1.9.6\n")

	lines := c.Run(t, "bud up", harness.Timeout(2*time.Minute))
	harness.OutputContains(t, lines, "activated: python[3.9.0]")

	lines = c.Run(t, "pip freeze")
	harness.OutputContains(t, lines, "pkginfo==1.9.6")
}
