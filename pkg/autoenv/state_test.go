package autoenv

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/stretchr/testify/require"
)

func newStateManager(env *env.Env) *StateManager {
	return &StateManager{env: env}
}

func TestStateSerialization(t *testing.T) {
	env := env.New([]string{})
	state := newStateManager(env)

	state.SetProjectSlug("p-1")
	require.Equal(t, `{"project":"p-1","features":null}`,
		env.Get("__BUD_AUTOENV"))

	state.SetFeature(&FeatureInfo{"f1", "v1"})
	require.Equal(t, `{"project":"p-1","features":[{"name":"f1","param":"v1"}]}`,
		env.Get("__BUD_AUTOENV"))

	state.SetProjectSlug("p-2")
	require.Equal(t, `{"project":"p-2","features":[{"name":"f1","param":"v1"}]}`,
		env.Get("__BUD_AUTOENV"))

	state.UnsetFeature("f1")
	require.Equal(t, `{"project":"p-2","features":[]}`,
		env.Get("__BUD_AUTOENV"))
}

func TestStateSetUnsetFeatures(t *testing.T) {
	env := env.New([]string{})

	newStateManager(env).SetFeature(NewFeatureInfo("rust", "v1"))
	require.Equal(t, "rust:v1", newStateManager(env).GetActiveFeatures().String())

	newStateManager(env).SetFeature(NewFeatureInfo("elixir", "v1"))
	require.Equal(t, "rust:v1 elixir:v1", newStateManager(env).GetActiveFeatures().String())

	newStateManager(env).SetFeature(NewFeatureInfo("rust", "v2"))
	require.Equal(t, "elixir:v1 rust:v2", newStateManager(env).GetActiveFeatures().String())

	newStateManager(env).UnsetFeature("elixir")
	require.Equal(t, "rust:v2", newStateManager(env).GetActiveFeatures().String())

	newStateManager(env).UnsetFeature("rust")
	require.Equal(t, "", newStateManager(env).GetActiveFeatures().String())
}

func TestStateSetGetProjectSlug(t *testing.T) {
	env := env.New([]string{})

	newStateManager(env).SetProjectSlug("p-1")
	require.Equal(t, "p-1", newStateManager(env).GetProjectSlug())

	newStateManager(env).SetProjectSlug("p-123")
	require.Equal(t, "p-123", newStateManager(env).GetProjectSlug())
}
