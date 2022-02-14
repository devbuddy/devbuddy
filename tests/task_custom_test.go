package integration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/tests/context"
)

var customTaskDevYml = `
up:
- custom:
    name: TestCustom
    met?: test -e sentinel
    meet: echo A > sentinel
`

func Test_Task_Custom(t *testing.T) {
	c := CreateContextAndInit(t)

	CreateProject(c, "project", customTaskDevYml)

	// file does not exist -> task must run
	c.Run("bud up")
	content := c.Cat("sentinel")
	require.Equal(t, "A", content)

	// file already exists -> task must not run
	c.Run("bud up")
	content = c.Cat("sentinel")
	require.Equal(t, "A", content) // same content
}

func Test_Task_Custom_Subdir(t *testing.T) {
	c := CreateContextAndInit(t)

	CreateProject(c, "project", customTaskDevYml)

	// The command must work in a sub-dir, but run in the project root
	c.Run("mkdir subdir")
	c.Run("cd subdir")
	c.Run("bud up")
	c.Run("cd ..")

	content := c.Cat("sentinel")
	require.Equal(t, "A", content)
}

func Test_Task_Custom_Fails(t *testing.T) {
	c := CreateContextAndInit(t)

	CreateProject(c, "project", `
up:
- custom:
    name: TestCustom
    met?: exit 1
    meet: exit 1
`)

	lines := c.Run("bud up", context.ExitCode(1))
	OutputContains(t, lines, "Running: exit 1", "The task action failed to run: command failed with exit code 1")
}

func Test_Task_Custom_With_Env_From_Shell(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Run("export MYVAR=poipoi")

	CreateProject(c, "project",
		`env:`,
		`  MYVAR: poipoi`,
		`up:`,
		`- custom:`,
		`    name: Test`,
		`    met?: echo $MYVAR > somefile`,
		`    meet: exit 0`,
	)
	c.Run("bud up")

	content := c.Cat("somefile")
	require.Equal(t, "poipoi", content)
}

func Test_Task_Custom_With_Env_At_First_Run(t *testing.T) {
	t.Skip("Fixme: env vars not set before tasks?")
	c := CreateContextAndInit(t)

	CreateProject(c, "project",
		`env:`,
		`  MYVAR: poipoi`,
		`up:`,
		`- custom:`,
		`    name: Test`,
		`    met?: echo $MYVAR > somefile`,
		`    meet: exit 0`,
	)

	c.Run("bud up")

	content := c.Cat("somefile")
	require.Equal(t, "poipoi", content)
}

func Test_Task_Custom_With_Env_Previously_Set_By_DevBuddy(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Run("export MYVAR=poipoi")

	p := CreateProject(c, "project",
		`env:`,
		`  MYVAR: poipoi`,
	)
	c.Run("bud up")

	p.UpdateDevYml(c,
		`env:`,
		`  MYVAR: poipoi`,
		`up:`,
		`- custom:`,
		`    name: Test`,
		`    met?: echo $MYVAR > somefile`,
		`    meet: exit 0`,
	)
	c.Run("bud up")

	content := c.Cat("somefile")
	require.Equal(t, "poipoi", content)
}
