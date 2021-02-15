package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Help(t *testing.T) {
	c := CreateContext(t)

	lines := c.Run("bud")
	OutputContains(t, lines, "Usage:", "Available Commands:", "Flags:")

	lines = c.Run("bud --help")
	OutputContains(t, lines, "Usage:", "Available Commands:", "Flags:")
}

func Test_Cmd_Version(t *testing.T) {
	c := CreateContext(t)

	lines := c.Run("bud --version")
	require.Equal(t, []string{"bud version devel"}, lines)
}

func Test_Cmd_DebugInfo(t *testing.T) {
	c := CreateContext(t)

	lines := c.Run("bud --debug-info")
	OutputContains(t, lines, "**DevBuddy version**", "`devel`")
	OutputContains(t, lines, "Project not found.")
}

func Test_Cmd_DebugInfo_Project(t *testing.T) {
	c := CreateContext(t)

	c.Run("mkdir poipoi")
	c.Cd("poipoi")
	c.Run("touch dev.yml")

	lines := c.Run("bud --debug-info")
	OutputContains(t, lines, "Project path: `/home/tester/poipoi`")
}
