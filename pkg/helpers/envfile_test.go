package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func TestLoadEnvfile(t *testing.T) {
	defer filet.CleanUp(t)
	envfilePath := filet.File(t, ".env", "\nPOIPOI=fromEnvFile")

	env := env.New([]string{})
	env.Set("EXISTING", "XXX")
	env.Set("POIPOI", "XXX")

	err := helpers.LoadEnvfile(env, envfilePath.Name())
	require.Nil(t, err)
	require.Equal(t, "fromEnvFile", env.Get("POIPOI"))
	require.Equal(t, "XXX", env.Get("EXISTING"))
}

func TestLoadEnvfileNoFile(t *testing.T) {
	defer filet.CleanUp(t)
	envfilePath := filet.File(t, "whatever", "")

	env := env.New([]string{})

	err := helpers.LoadEnvfile(env, envfilePath.Name()+"nope")
	require.NotNil(t, err)
	require.Zero(t, len(env.Environ()))
}
