package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Env_EnvFile(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Write("dev.yml", `up: [envfile]`)
	c.Write(".env", `TESTVAR=FooBAr`)

	lines := c.Run("bud up")
	OutputEqual(t, lines, "◼︎ EnvFile")

	value := c.GetEnv("TESTVAR")
	require.Equal(t, "FooBAr", value)

	// Clean the env when leaving the project directory
	c.Run("cd /")

	value = c.GetEnv("TESTVAR")
	require.Equal(t, "", value)
}

func Test_Env_EnvFile_In_Process(t *testing.T) {
	t.Skip("to be fixed")

	c := CreateContextAndInit(t)

	devYml := `
up:
- envfile
- custom:
    name: succeed if TESTVAR is set
    met?: test -n "${TESTVAR}"
    meet: echo "TESTVAR is not set"; false
`
	c.Write("dev.yml", devYml)
	c.Write(".env", `TESTVAR=FooBAr`)

	c.Run("bud up")
}

func Test_Env_EnvFile_Changes(t *testing.T) {
	t.Skip("to be fixed")

	c := CreateContextAndInit(t)

	c.Write("dev.yml", `up: [envfile]`)
	c.Write(".env", `TESTVAR=one`)
	c.Run("bud up")

	value := c.GetEnv("TESTVAR")
	require.Equal(t, "one", value)

	// Change .env file
	c.Write(".env", `TESTVAR=two`)
	c.Run("bud up")

	value = c.GetEnv("TESTVAR")
	require.Equal(t, "two", value)
}
