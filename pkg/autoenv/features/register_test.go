package features

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/context"

	"github.com/stretchr/testify/require"
)

type testFeature string

func (t testFeature) Name() string {
	return string(t)
}

func (testFeature) Activate(ctx *context.Context, param string) (bool, error) {
	return false, nil
}

func (testFeature) Deactivate(ctx *context.Context, param string) {}

func TestRegister(t *testing.T) {
	reg := Register{}

	reg.Register(python{})
	reg.Register(golang{})

	require.ElementsMatch(t, reg.Names(), []string{"python", "golang"})

	env := reg.Get("python")
	require.Equal(t, env.Name(), "python")

	env = reg.Get("golang")
	require.Equal(t, env.Name(), "golang")

	env = reg.Get("nope")
	require.Nil(t, env, "unknown feature: nope")
}
