package features

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/context"
	envpkg "github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/termui"

	"github.com/stretchr/testify/require"
)

func newTestContext() *context.Context {
	_, ui := termui.NewTesting(false)
	return &context.Context{
		Env: envpkg.New([]string{}),
		UI:  ui,
	}
}

func TestEnvActivateWithJSONParam(t *testing.T) {
	ctx := newTestContext()

	devUpNeeded, err := env{}.Activate(ctx, `{"FOO":"bar","BAZ":"qux"}`)
	require.NoError(t, err)
	require.False(t, devUpNeeded)

	require.Equal(t, "bar", ctx.Env.Get("FOO"))
	require.Equal(t, "qux", ctx.Env.Get("BAZ"))
}

func TestEnvActivateWithInvalidJSON(t *testing.T) {
	ctx := newTestContext()

	_, err := env{}.Activate(ctx, "not-json")
	require.Error(t, err)
	require.Contains(t, err.Error(), "env feature: invalid param")
}

func TestEnvActivateWithEmptyObject(t *testing.T) {
	ctx := newTestContext()

	devUpNeeded, err := env{}.Activate(ctx, `{}`)
	require.NoError(t, err)
	require.False(t, devUpNeeded)
}
