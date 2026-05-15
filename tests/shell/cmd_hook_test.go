package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Hook_Preserves_Previous_Exit_Code(t *testing.T) {
	c := CreateContextAndInit(t)

	p := CreateProject(t, c,
		`up: []`,
	)
	c.Cd(t, p.Path)

	lines := c.Run(t, "false; __bud_prompt_command; echo $?")
	OutputEqual(t, lines, "1")
}

func Test_Hook_DeferInit_Skips_First_Invocation(t *testing.T) {
	c := CreateContext(t)

	// Create a project and cd into it BEFORE shell init.
	// No PROMPT_COMMAND is set yet, so no hook fires during these steps.
	p := CreateProject(t, c,
		`env:`,
		`  TESTVAR: hello`,
	)
	c.Cd(t, p.Path)

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
	c := CreateContext(t)

	// Create a project and cd into it BEFORE shell init.
	p := CreateProject(t, c,
		`env:`,
		`  TESTVAR: hello`,
	)
	c.Cd(t, p.Path)

	// Source the init but unset PROMPT_COMMAND/precmd so we control hook calls.
	c.Run(t, `eval "$(bud --shell-init)"; unset PROMPT_COMMAND; precmd_functions=()`)

	// First hook call without defer should activate immediately.
	c.Run(t, `__bud_prompt_command`)
	value := c.GetEnv(t, "TESTVAR")
	require.Equal(t, "hello", value, "TESTVAR should be set after first hook")
}
