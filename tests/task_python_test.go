package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Env_Python(t *testing.T) {
	c := CreateContextAndInit(t)

	devYml := `
up:
- python: 3.9.0
`
	CreateProject(t, c, "project", devYml)

	lines := c.Run("bud up", context.Timeout(2*time.Minute))
	OutputContains(t, lines, "python activated. (3.9.0)")

	lines = c.Run("python --version")
	OutputEqual(t, lines, "Python 3.9.0")

	// Assert that the virtualenv is active
	lines = c.Run("python -c 'import sys; print(sys.prefix)'")
	OutputContains(t, lines, "/.local/share/bud/virtualenvs/")
}
