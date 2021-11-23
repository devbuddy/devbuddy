package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/test"
)

func TestLoadEnvfile(t *testing.T) {
	_, tmpfile := test.CreateFile(t, ".env", []byte("\nPOIPOI=fromEnvFile"))

	env := env.New([]string{})
	env.Set("EXISTING", "XXX")
	env.Set("POIPOI", "XXX")

	err := helpers.LoadEnvfile(env, tmpfile)
	require.Nil(t, err)
	require.Equal(t, "fromEnvFile", env.Get("POIPOI"))
	require.Equal(t, "XXX", env.Get("EXISTING"))
}

func TestLoadEnvfileNoFile(t *testing.T) {
	tmpdir := t.TempDir()

	env := env.New([]string{})

	err := helpers.LoadEnvfile(env, tmpdir+"/nope")
	require.NotNil(t, err)
	require.Zero(t, len(env.Environ()))
}
