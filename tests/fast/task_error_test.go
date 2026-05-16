package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

func Test_Task_Error_NotAList(t *testing.T) {
	c := harness.NewCLI(t)

	c.Write(t, "dev.yml", `up: somestring`)

	lines := c.Run(t, "bud up", harness.ExitCode(1))
	harness.OutputContains(t, lines, "string was used where sequence is expected")
}

func Test_Task_Error_InvalidType(t *testing.T) {
	c := harness.NewCLI(t)

	c.Write(t, "dev.yml", `up: [true]`)

	lines := c.Run(t, "bud up", harness.ExitCode(1))
	harness.OutputEqual(t, lines, `Error: parsing task: invalid task: "true"`)
}

func Test_Task_Error_UnknownTask(t *testing.T) {
	c := harness.NewCLI(t)

	c.Write(t, "dev.yml", `up: [notatask]`)

	lines := c.Run(t, "bud up", harness.ExitCode(1))
	harness.OutputContains(t, lines, `unknown task: "notatask"`)
}

func Test_Task_Error_Invalid_Hash_Type(t *testing.T) {
	c := harness.NewCLI(t)

	c.Write(t, "dev.yml", `up: [ { go: {version: 1.16} } ]`)

	lines := c.Run(t, "bud up", harness.ExitCode(1))
	harness.OutputEqual(t, lines,
		`Error: task "go": key "version": expecting a string, found a float64 (1.16)`,
	)
}

func Test_Task_Error_Invalid_List(t *testing.T) {
	c := harness.NewCLI(t)

	c.Write(t, "dev.yml", `up: [ { homebrew: {} } ]`)

	lines := c.Run(t, "bud up", harness.ExitCode(1))
	harness.OutputEqual(t, lines,
		`Error: task "homebrew": expecting a list of strings, found a map[string]interface {} (map[])`,
	)
}
