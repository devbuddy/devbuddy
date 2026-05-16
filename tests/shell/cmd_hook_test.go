package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
	"github.com/stretchr/testify/require"
)

func Test_Hook_Preserves_Previous_Exit_Code(t *testing.T) {
	c := harness.NewDockerInit(t)

	harness.NewDockerProject(t, c,
		`up: []`,
	)

	lines := c.Run(t, "false; __bud_prompt_command; echo $?")
	harness.OutputEqual(t, lines, "1")
}

func Test_Hook_DeferInit_Skips_First_Invocation(t *testing.T) {
	c := harness.NewDocker(t)

	// Create a project and cd into it BEFORE shell init.
	// No PROMPT_COMMAND is set yet, so no hook fires during these steps.
	harness.NewDockerProject(t, c,
		`env:`,
		`  TESTVAR: hello`,
	)

	// Source the init but unset PROMPT_COMMAND/precmd so we control hook calls manually.
	c.Run(t, `__bud_defer_init=1; eval "$(bud --shell-init)"; unset PROMPT_COMMAND; precmd_functions=()`)

	// Call the hook manually for the first time — should be deferred.
	c.Run(t, `__bud_prompt_command`)
	value := c.GetEnv(t, "TESTVAR")
	require.Equal(t, "", value, "TESTVAR should not be set after deferred first hook")

	// Call the hook again — should activate now.
	c.Run(t, `__bud_prompt_command`)
	value = c.GetEnv(t, "TESTVAR")
	require.Equal(t, "hello", value, "TESTVAR should be set after second hook")
}

func Test_Hook_Without_DeferInit_Activates_Immediately(t *testing.T) {
	c := harness.NewDocker(t)

	// Create a project and cd into it BEFORE shell init.
	harness.NewDockerProject(t, c,
		`env:`,
		`  TESTVAR: hello`,
	)

	// Source the init but unset PROMPT_COMMAND/precmd so we control hook calls.
	c.Run(t, `eval "$(bud --shell-init)"; unset PROMPT_COMMAND; precmd_functions=()`)

	// First hook call without defer should activate immediately.
	c.Run(t, `__bud_prompt_command`)
	value := c.GetEnv(t, "TESTVAR")
	require.Equal(t, "hello", value, "TESTVAR should be set after first hook")
}
