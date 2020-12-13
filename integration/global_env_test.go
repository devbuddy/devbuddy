package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/integration/context"
	"github.com/stretchr/testify/require"
)

func Test_Env(t *testing.T) {
	c := CreateContextAndInit(t)

	devYml := `env: {TESTVAR: TESTVALUE}`

	c.Write("dev.yml", devYml)

	lines := c.Run("bud up", context.ExitCode(0))
	OutputEqual(t, lines, "◼︎ Env")

	value := c.GetEnv("TESTVAR")
	require.Equal(t, "TESTVALUE", value)

	// Clean the env when leaving the project directory
	c.Run("cd /")

	value = c.GetEnv("TESTVAR")
	require.Equal(t, "", value)
}
