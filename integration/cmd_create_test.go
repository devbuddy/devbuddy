package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Create(t *testing.T) {
	c := CreateContextAndInit(t)

	output := c.Run("bud create orgname/projname")
	require.Equal(t, []string{
		"🐼  Creating a default dev.yml file.",
		"⚠️   Open dev.yml to adjust for your needs.",
		"🐼  jumping to github.com:orgname/projname",
		"🐼  env activated.",
	}, output)

	cwd := c.Cwd()
	require.Equal(t, "/home/tester/src/github.com/orgname/projname", cwd)

	files := c.Ls(".")
	require.ElementsMatch(t, []string{"dev.yml"}, files)

	devFile := c.Cat("dev.yml")
	require.Contains(t, devFile, "# DevBuddy config file")
}

func Test_Cmd_Create_Already_Exists(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Run("mkdir -p /home/tester/src/github.com/orgname/projname")

	lines := c.Run("bud create orgname/projname")
	OutputContains(t, lines, "project already exists locally")

	cwd := c.Cwd()
	require.Equal(t, "/home/tester/src/github.com/orgname/projname", cwd)
}
