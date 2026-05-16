package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
	"github.com/stretchr/testify/require"
)

const (
	devYmlMyCmd = `
commands:
  mycmd:
    run: touch somefile
`
	devYmlMyCmdShort = `
commands:
  mycmd: touch somefile
`
)

func Test_Cmd_Custom(t *testing.T) {
	c := harness.NewCLI(t)

	harness.NewCLIProject(t, c, devYmlMyCmd)

	lines := c.Run(t, "bud mycmd")
	harness.OutputEqual(t, lines, "🐼  running touch somefile")

	files := c.Ls(t, ".")
	require.ElementsMatch(t, files, []string{"dev.yml", "somefile"})
}

func Test_Cmd_Custom_Short_Syntax(t *testing.T) {
	c := harness.NewCLI(t)

	harness.NewCLIProject(t, c, devYmlMyCmdShort)

	lines := c.Run(t, "bud mycmd")
	harness.OutputEqual(t, lines, "🐼  running touch somefile")
}

func Test_Cmd_Custom_Envs_Are_Applied(t *testing.T) {
	c := harness.NewCLI(t)

	harness.NewCLIProject(t, c,
		`env:`,
		`  MYVAR: poipoi`,
		`commands:`,
		`  mycmd: "echo __${MYVAR}__ > result"`,
	)

	c.Run(t, "bud mycmd")
	c.AssertContains(t, "result", "__poipoi__")
}

func Test_Cmd_Custom_With_Piped_Stdin(t *testing.T) {
	c := harness.NewCLI(t)

	harness.NewCLIProject(t, c, devYmlMyCmd)

	lines := c.Run(t, "echo '' | bud mycmd")
	harness.OutputEqual(t, lines, "🐼  running touch somefile")
}

func Test_Cmd_Custom_Output(t *testing.T) {
	c := harness.NewCLI(t)

	harness.NewCLIProject(t, c,
		`commands:`,
		`  greet: echo "hello world"`,
	)

	lines := c.Run(t, "bud greet")
	harness.OutputContains(t, lines, "hello world")
}

func Test_Cmd_Custom_Exit_Code(t *testing.T) {
	c := harness.NewCLI(t)

	harness.NewCLIProject(t, c,
		`commands:`,
		`  fail: exit 42`,
	)

	// bud reports the failure but normalizes the exit code to 1
	lines := c.Run(t, "bud fail", harness.ExitCode(1))
	harness.OutputContains(t, lines, "command failed with exit code 42")
}

func Test_Cmd_Custom_Always_Run_In_Project_Root(t *testing.T) {
	c := harness.NewCLI(t)

	harness.NewCLIProject(t, c, devYmlMyCmd)
	c.Run(t, "mkdir foobar")
	c.Cd(t, "foobar")

	lines := c.Run(t, "bud mycmd")
	harness.OutputEqual(t, lines, "🐼  running touch somefile")

	files := c.Ls(t, "..")
	require.ElementsMatch(t, files, []string{"dev.yml", "foobar", "somefile"})
}
