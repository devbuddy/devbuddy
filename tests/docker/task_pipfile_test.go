package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/context"
	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

func Test_Task_Pipfile(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	harness.NewDockerProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- pipfile`,
	)

	c.Write(t, "Pipfile", "[packages]\n\"pkginfo\" = \"==1.9.6\"\n")

	c.Run(t, "bud up", context.Timeout(2*time.Minute))

	lines := c.Run(t, "pip freeze")
	harness.OutputContains(t, lines, "pkginfo==1.9.6")
}
