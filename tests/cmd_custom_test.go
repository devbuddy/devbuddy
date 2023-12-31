package integration

import (
	"testing"

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
	c := CreateContextAndInit(t)

	project := CreateProject(t, c, devYmlMyCmd)
	c.Cd(t, project.Path)

	lines := c.Run(t, "bud mycmd")
	OutputEqual(t, lines, "🐼  running touch somefile")

	files := c.Ls(t, ".")
	require.ElementsMatch(t, files, []string{"dev.yml", "somefile"})
}

func Test_Cmd_Custom_Short_Syntax(t *testing.T) {
	c := CreateContextAndInit(t)

	project := CreateProject(t, c, devYmlMyCmdShort)
	c.Cd(t, project.Path)

	lines := c.Run(t, "bud mycmd")
	OutputEqual(t, lines, "🐼  running touch somefile")
}

func Test_Cmd_Custom_Envs_Are_Applied(t *testing.T) {
	c := CreateContextAndInit(t)

	project := CreateProject(t, c,
		`env:`,
		`  MYVAR: poipoi`,
		`commands:`,
		`  mycmd: "echo __${MYVAR}__ > result"`,
	)
	c.Cd(t, project.Path)

	c.Run(t, "bud mycmd")
	c.AssertContains(t, "result", "__poipoi__")
}

func Test_Cmd_Custom_Always_Run_In_Project_Root(t *testing.T) {
	c := CreateContextAndInit(t)

	project := CreateProject(t, c, devYmlMyCmd)
	c.Cd(t, project.Path)
	c.Run(t, "mkdir foobar")
	c.Cd(t, "foobar")

	lines := c.Run(t, "bud mycmd")
	OutputEqual(t, lines, "🐼  running touch somefile")

	files := c.Ls(t, "..")
	require.ElementsMatch(t, files, []string{"dev.yml", "foobar", "somefile"})
}
