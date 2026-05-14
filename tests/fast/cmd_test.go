package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Help(t *testing.T) {
	c := CreateContext(t)

	lines := c.Run(t, "bud")
	OutputContains(t, lines, "Usage:", "DevBuddy Commands:", "  cd", "  clone", "  up", "Flags:")

	lines = c.Run(t, "bud --help")
	OutputContains(t, lines, "Usage:", "DevBuddy Commands:", "  cd", "  clone", "  up", "Flags:")
}

func Test_Cmd_Version(t *testing.T) {
	c := CreateContext(t)

	lines := c.Run(t, "bud --version")
	require.Equal(t, []string{"bud version devel"}, lines)
}

func Test_Cmd_DebugInfo(t *testing.T) {
	c := CreateContext(t)

	lines := c.Run(t, "bud --debug-info")
	OutputContains(t, lines, "**DevBuddy version**", "`devel`")
	OutputContains(t, lines, "Project not found.")
}

func Test_Cmd_DebugInfo_Project(t *testing.T) {
	c := CreateContext(t)

	c.Run(t, "mkdir poipoi")
	c.Cd(t, "poipoi")
	c.Run(t, "touch dev.yml")

	lines := c.Run(t, "bud --debug-info")
	OutputContains(t, lines, "Project path: `/home/tester/poipoi`")
}
