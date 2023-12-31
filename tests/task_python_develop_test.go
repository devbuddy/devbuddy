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

	devYml := `
up:
- python: 3.9.0
- python_develop
`
	CreateProject(t, c, devYml)

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

	devYml := `
up:
- python: 3.9.0
- python_develop:
    extras: [test]
`
	CreateProject(t, c, devYml)

	c.Write(t, "setup.py", generateTestSetupPy(1))

	c.Run(t, "bud up", context.Timeout(2*time.Minute))

	lines := c.Run(t, "pip freeze")
	OutputContains(t, lines, "pyreleaser==0.5.2")
}

func Test_Task_Python_Develop_Without_Extra_Packages(t *testing.T) {
	c := CreateContextAndInit(t)

	devYml := `
up:
- python: 3.9.0
- python_develop:
`
	CreateProject(t, c, devYml)

	c.Write(t, "setup.py", generateTestSetupPy(1))

	c.Run(t, "bud up", context.Timeout(2*time.Minute))

	lines := c.Run(t, "pip freeze")
	OutputNotContain(t, lines, "pyreleaser==0.5.2")
}

func generateTestSetupPy(version int) string {
	lines := []string{
		`from setuptools import setup, find_packages`,
		fmt.Sprintf(`setup(name='devbuddy-test-pkg', version='%d', extras_require={'test': ['pyreleaser==0.5.2']})`, version),
		`open("sentinel", "w").write("")`,
	}
	return strings.Join(lines, "\n")
}
