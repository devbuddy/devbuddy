package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

func Test_Task_Python_Develop(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	harness.NewProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- python_develop`,
	)

	// Install in develop mode

	c.Write(t, "setup.py", generateTestSetupPy(42))

	lines := c.Run(t, "bud up", harness.Timeout(2*time.Minute))
	harness.OutputContains(t, lines, "activated: python[3.9.0]")

	lines = c.Run(t, "pip show devbuddy-test-pkg")
	harness.OutputContains(t, lines, "Version: 42")

	// Update the package

	c.Write(t, "setup.py", generateTestSetupPy(84))

	c.Run(t, "bud up", harness.Timeout(2*time.Minute))

	lines = c.Run(t, "pip show devbuddy-test-pkg")
	harness.OutputContains(t, lines, "Version: 84")
}

func Test_Task_Python_Develop_With_Extra_Packages(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	harness.NewProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- python_develop:`,
		`    extras: [test]`,
	)

	c.Write(t, "setup.py", generateTestSetupPy(1))

	c.Run(t, "bud up", harness.Timeout(2*time.Minute))

	lines := c.Run(t, "pip freeze")
	harness.OutputContains(t, lines, "pkginfo==1.9.6")
}

func Test_Task_Python_Develop_Without_Extra_Packages(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	harness.NewProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- python_develop:`,
	)

	c.Write(t, "setup.py", generateTestSetupPy(1))

	c.Run(t, "bud up", harness.Timeout(2*time.Minute))

	lines := c.Run(t, "pip freeze")
	harness.OutputNotContains(t, lines, "pkginfo==1.9.6")
}

func generateTestSetupPy(version int) string {
	lines := []string{
		`from setuptools import setup, find_packages`,
		fmt.Sprintf(`setup(name='devbuddy-test-pkg', version='%d', extras_require={'test': ['pkginfo==1.9.6']})`, version),
	}
	return strings.Join(lines, "\n")
}
