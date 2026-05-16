package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/context"
	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

func Test_Env_Python(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	harness.NewDockerProject(t, c,
		`up:`,
		`- python: 3.9.0`,
	)

	lines := c.Run(t, "bud up", context.Timeout(2*time.Minute))
	harness.OutputContains(t, lines, "activated: python[3.9.0]")

	lines = c.Run(t, "python --version")
	harness.OutputEqual(t, lines, "Python 3.9.0")

	// Assert that the virtualenv is active
	lines = c.Run(t, "python -c 'import sys; print(sys.prefix)'")
	harness.OutputContains(t, lines, "/.local/share/bud/virtualenvs/")
}
