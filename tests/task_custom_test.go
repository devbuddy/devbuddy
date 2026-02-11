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
	c, _ := CreateContextAndProject(t, customTaskDevYml)

	// file does not exist -> task must run
	c.Run(t, "bud up")
	content := c.Cat(t, "sentinel")
	require.Equal(t, "A", content)

	// file already exists -> task must not run
	c.Run(t, "bud up")
	content = c.Cat(t, "sentinel")
	require.Equal(t, "A", content) // same content
}

func Test_Task_Custom_Subdir(t *testing.T) {
	c, _ := CreateContextAndProject(t, customTaskDevYml)

	// The command must work in a sub-dir, but run in the project root
	c.Run(t, "mkdir subdir")
	c.Run(t, "cd subdir")
	c.Run(t, "bud up")
	c.Run(t, "cd ..")

	content := c.Cat(t, "sentinel")
	require.Equal(t, "A", content)
}

func Test_Task_Custom_Fails(t *testing.T) {
	c, _ := CreateContextAndProject(t,
		`up:`,
		`- custom:`,
		`    name: TestCustom`,
		`    met?: exit 1`,
		`    meet: exit 1`,
	)

	lines := c.Run(t, "bud up", context.ExitCode(1))
	OutputContains(t, lines, "Running: exit 1", `action "": failed to run: command failed with exit code 1`)
}

func Test_Task_Custom_With_Env_From_Shell(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Run(t, "export MYVAR=poipoi")

	p := CreateProject(t, c,
		`env:`,
		`  MYVAR: poipoi`,
		`up:`,
		`- custom:`,
		`    name: Test`,
		`    met?: echo $MYVAR > somefile`,
		`    meet: exit 0`,
	)
	c.Cd(t, p.Path)

	c.Run(t, "bud up")

	content := c.Cat(t, "somefile")
	require.Equal(t, "poipoi", content)
}

func Test_Task_Custom_With_Env_At_First_Run(t *testing.T) {
	c, _ := CreateContextAndProject(t,
		`env:`,
		`  MYVAR: poipoi`,
		`up:`,
		`- custom:`,
		`    name: Test`,
		`    met?: echo $MYVAR > somefile`,
		`    meet: exit 0`,
	)

	c.Run(t, "bud up")

	content := c.Cat(t, "somefile")
	require.Equal(t, "poipoi", content)
}

func Test_Task_Custom_With_Env_Previously_Set_By_DevBuddy(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Run(t, "export MYVAR=poipoi")

	p := CreateProject(t, c,
		`env:`,
		`  MYVAR: poipoi`,
	)
	c.Cd(t, p.Path)

	c.Run(t, "bud up")

	p.WriteDevYml(t,
		`env:`,
		`  MYVAR: poipoi`,
		`up:`,
		`- custom:`,
		`    name: Test`,
		`    met?: echo $MYVAR > somefile`,
		`    meet: exit 0`,
	)
	c.Run(t, "bud up")

	content := c.Cat(t, "somefile")
	require.Equal(t, "poipoi", content)
}
