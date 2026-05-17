package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Help(t *testing.T) {
	c := harness.NewCLI(t)

	lines := c.Run(t, "bud")
	harness.OutputContains(t, lines, "Usage:", "DevBuddy Commands:", "  cd", "  clone", "  up", "Flags:")

	lines = c.Run(t, "bud --help")
	harness.OutputContains(t, lines, "Usage:", "DevBuddy Commands:", "  cd", "  clone", "  up", "Flags:")
}

func Test_Cmd_Version(t *testing.T) {
	c := harness.NewCLI(t)

	lines := c.Run(t, "bud --version")
	require.Equal(t, []string{"bud version devel"}, lines)
}

func Test_Cmd_UnknownCommand_ShowsHelpHint(t *testing.T) {
	c := harness.NewCLI(t)

	lines := c.Run(t, "bud does-not-exist", harness.ExitCode(1))

	harness.OutputContains(t, lines,
		`Error: unknown command "does-not-exist" for "bud"`,
		`Run 'bud --help' for usage.`,
	)
}

func Test_Cmd_UnknownNestedCommand_ShowsNestedHelpHint(t *testing.T) {
	c := harness.NewCLI(t)

	lines := c.Run(t, "bud tree does-not-exist", harness.ExitCode(1))

	harness.OutputContains(t, lines,
		`Error: unknown command "does-not-exist" for "bud tree"`,
		`Run 'bud tree --help' for usage.`,
	)
}

func Test_Cmd_DebugInfo(t *testing.T) {
	c := harness.NewCLI(t)

	lines := c.Run(t, "bud --debug-info")
	harness.OutputContains(t, lines, "**DevBuddy version**", "`devel`")
	harness.OutputContains(t, lines, "Project not found.")
}

func Test_Cmd_DebugInfo_Project(t *testing.T) {
	c := harness.NewCLI(t)

	c.Run(t, "mkdir poipoi")
	c.Cd(t, "poipoi")
	c.Run(t, "touch dev.yml")

	lines := c.Run(t, "bud --debug-info")
	harness.OutputContains(t, lines, "Project path: `"+c.Cwd(t)+"`")
}
