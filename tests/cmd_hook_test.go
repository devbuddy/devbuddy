package integration

import (
	"testing"
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
