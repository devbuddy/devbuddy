package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/integration/context"
)

func Test_Cmd_Up_Invalid(t *testing.T) {
	c := CreateContextAndInit(t)

	devYml := `up: somestring`

	c.Write("dev.yml", devYml)

	lines := c.Run("bud up", context.ExitCode(1))
	OutputEqual(t, lines,
		"Error: yaml: unmarshal errors:",
		"  line 1: cannot unmarshal !!str `somestring` into []interface {}",
	)

	devYml = `up: [notatask]`

	c.Write("dev.yml", devYml)

	lines = c.Run("bud up", context.ExitCode(0)) // TODO: This should probably return 1
	OutputEqual(t, lines,
		"◼︎ Unknown",
		`  Warning: Unknown task: "notatask"`,
	)

	devYml = `up: [true]`

	c.Write("dev.yml", devYml)

	lines = c.Run("bud up", context.ExitCode(1))
	OutputEqual(t, lines, `Error: invalid task: "true"`)

}
