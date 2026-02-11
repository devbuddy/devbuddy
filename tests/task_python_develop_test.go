package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Task_Python_Develop(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- python_develop`,
	)
	c.Cd(t, p.Path)

	// Install in develop mode

	c.Write(t, "setup.py", generateTestSetupPy(42))

	lines := c.Run(t, "bud up", context.Timeout(2*time.Minute))
	OutputContains(t, lines, "python activated. (3.9.0)")

	lines = c.Run(t, "pip show devbuddy-test-pkg")
	OutputContains(t, lines, "Version: 42")

	// Update the package

	c.Write(t, "setup.py", generateTestSetupPy(84))

	c.Run(t, "bud up", context.Timeout(2*time.Minute))

	lines = c.Run(t, "pip show devbuddy-test-pkg")
	OutputContains(t, lines, "Version: 84")
}

func Test_Task_Python_Develop_With_Extra_Packages(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- python_develop:`,
		`    extras: [test]`,
	)
	c.Cd(t, p.Path)

	c.Write(t, "setup.py", generateTestSetupPy(1))

	c.Run(t, "bud up", context.Timeout(2*time.Minute))

	lines := c.Run(t, "pip freeze")
	OutputContains(t, lines, "pkginfo==1.9.6")
}

func Test_Task_Python_Develop_Without_Extra_Packages(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up:`,
		`- python: 3.9.0`,
		`- python_develop:`,
	)
	c.Cd(t, p.Path)

	c.Write(t, "setup.py", generateTestSetupPy(1))

	c.Run(t, "bud up", context.Timeout(2*time.Minute))

	lines := c.Run(t, "pip freeze")
	OutputNotContains(t, lines, "pkginfo==1.9.6")
}

func generateTestSetupPy(version int) string {
	lines := []string{
		`from setuptools import setup, find_packages`,
		fmt.Sprintf(`setup(name='devbuddy-test-pkg', version='%d', extras_require={'test': ['pkginfo==1.9.6']})`, version),
	}
	return strings.Join(lines, "\n")
}
