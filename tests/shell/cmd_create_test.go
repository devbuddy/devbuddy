package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Create(t *testing.T) {
	c := CreatePTYContextAndInit(t)

	output := c.Run(t, "bud create --template default orgname/projname")
	require.Equal(t, []string{
		"🐼  Created dev.yml with template default",
		"⚠️   Open dev.yml to adjust for your needs.",
		"🐼  jumping to github.com:orgname/projname",
		"🐼  activated: env",
	}, output)

	cwd := c.Cwd(t)
	require.Equal(t, "/home/tester/src/github.com/orgname/projname", cwd)

	files := c.Ls(t, ".")
	require.ElementsMatch(t, []string{"dev.yml"}, files)

	devFile := c.Cat(t, "dev.yml")
	require.Contains(t, devFile, "# DevBuddy config file")
}

func Test_Cmd_Create_Already_Exists(t *testing.T) {
	c := CreatePTYContextAndInit(t)

	c.Run(t, "mkdir -p /home/tester/src/github.com/orgname/projname")

	lines := c.Run(t, "bud create orgname/projname")
	OutputContains(t, lines, "project already exists locally")

	cwd := c.Cwd(t)
	require.Equal(t, "/home/tester/src/github.com/orgname/projname", cwd)
}
