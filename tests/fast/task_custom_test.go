package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
	"github.com/stretchr/testify/require"
)

var customTaskDevYml = `
up:
- custom:
    name: TestCustom
    met?: test -e sentinel
    meet: echo A > sentinel
`

func Test_Task_Custom(t *testing.T) {
	c := harness.NewCLI(t)
	harness.NewCLIProject(t, c, customTaskDevYml)

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
	c := harness.NewCLI(t)
	harness.NewCLIProject(t, c, customTaskDevYml)

	// The command must work in a sub-dir, but run in the project root
	c.Run(t, "mkdir subdir")
	c.Cd(t, "subdir")
	c.Run(t, "bud up")
	c.Cd(t, "..")

	content := c.Cat(t, "sentinel")
	require.Equal(t, "A", content)
}

func Test_Task_Custom_Fails(t *testing.T) {
	c := harness.NewCLI(t)
	harness.NewCLIProject(t, c,
		`up:`,
		`- custom:`,
		`    name: TestCustom`,
		`    met?: exit 1`,
		`    meet: exit 1`,
	)

	lines := c.Run(t, "bud up", harness.ExitCode(1))
	harness.OutputContains(t, lines, "Running: exit 1", `action "": failed to run: command failed with exit code 1`)
}

func Test_Task_Custom_With_Env_From_Shell(t *testing.T) {
	c := harness.NewCLI(t)

	c.Setenv("MYVAR", "poipoi")

	harness.NewCLIProject(t, c,
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

func Test_Task_Custom_With_Env_At_First_Run(t *testing.T) {
	c := harness.NewCLI(t)
	harness.NewCLIProject(t, c,
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
	c := harness.NewCLI(t)

	c.Setenv("MYVAR", "poipoi")

	harness.NewCLIProject(t, c,
		`env:`,
		`  MYVAR: poipoi`,
	)

	c.Run(t, "bud up")

	c.WriteLines(t, "dev.yml",
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
