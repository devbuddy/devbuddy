package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/context"
	"github.com/stretchr/testify/require"
)

func Test_Env(t *testing.T) {
	c := CreateContextAndInit(t)

	devYml := `env: {TESTVAR: TESTVALUE}`

	c.Write(t, "dev.yml", devYml)

	lines := c.Run(t, "bud up", context.ExitCode(0))
	OutputEqual(t, lines, "◼︎ Env")

	value := c.GetEnv(t, "TESTVAR")
	require.Equal(t, "TESTVALUE", value)

	// Clean the env when leaving the project directory
	c.Run(t, "cd /")

	value = c.GetEnv(t, "TESTVAR")
	require.Equal(t, "", value)
}
