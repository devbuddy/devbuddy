package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Env_Python(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up:`,
		`- python: 3.9.0`,
	)
	c.Cd(t, p.Path)

	lines := c.Run(t, "bud up", context.Timeout(2*time.Minute))
	OutputContains(t, lines, "python activated. (3.9.0)")

	lines = c.Run(t, "python --version")
	OutputEqual(t, lines, "Python 3.9.0")

	// Assert that the virtualenv is active
	lines = c.Run(t, "python -c 'import sys; print(sys.prefix)'")
	OutputContains(t, lines, "/.local/share/bud/virtualenvs/")
}
