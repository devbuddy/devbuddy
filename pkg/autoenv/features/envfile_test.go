package features

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/test"
)

func TestEnvfile(t *testing.T) {
	tmpdir, _ := test.CreateFile(t, ".env", []byte("\nPOIPOI=fromEnvFile"))
	os.Chdir(tmpdir)

	env := env.New([]string{})
	env.Set("EXISTING", "XXX")
	env.Set("POIPOI", "XXX")

	ctx := &context.Context{Env: env}

	feature := Envfile{}
	devneeded, err := feature.Activate(ctx, ".env")
	require.Nil(t, err)
	require.False(t, devneeded)

	require.Equal(t, "fromEnvFile", env.Get("POIPOI"))
	require.Equal(t, "XXX", env.Get("EXISTING"))
}

func TestEnvfileNoFile(t *testing.T) {
	tmpdir := t.TempDir()
	os.Chdir(tmpdir)

	env := env.New([]string{})

	ctx := &context.Context{Env: env}

	feature := Envfile{}
	devneeded, err := feature.Activate(ctx, ".env")
	require.EqualError(t, err, "open .env: no such file or directory")
	require.True(t, devneeded)

	require.Zero(t, len(env.Environ()))
}
