package autoenv

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/termui"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/stretchr/testify/require"
)

func newStateManager(env *env.Env) *StateManager {
	_, ui := termui.NewTesting(false)
	return &StateManager{env: env, UI: ui}
}

func TestStateSerialization(t *testing.T) {
	env := env.New([]string{})
	state := newStateManager(env)

	state.SetProjectSlug("p-1")
	require.Equal(t, `{"project":"p-1","features":null,"saved_env":{}}`,
		env.Get("__BUD_AUTOENV"))

	state.SetFeature(&FeatureInfo{"f1", "v1"})
	require.Equal(t, `{"project":"p-1","features":[{"name":"f1","param":"v1"}],"saved_env":{}}`,
		env.Get("__BUD_AUTOENV"))

	state.SetProjectSlug("p-2")
	require.Equal(t, `{"project":"p-2","features":[{"name":"f1","param":"v1"}],"saved_env":{}}`,
		env.Get("__BUD_AUTOENV"))

	state.UnsetFeature("f1")
	require.Equal(t, `{"project":"p-2","features":[],"saved_env":{}}`,
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

func TestStateSavedEnvDirect(t *testing.T) {
	env1 := env.New([]string{"GO111MODULES=off"})

	env1.Set("GOROOT", "/go/1") // imaginary project
	env1.Set("GO111MODULES", "on")
	newStateManager(env1).SaveEnv()

	env2 := env.New(env1.Environ()) // next prompt

	newStateManager(env2).RestoreEnv()

	require.False(t, env2.Has("GOROOT"))
	require.Equal(t, "off", env2.Get("GO111MODULES"))
}

func TestStateSavedEnvMultipleChanges(t *testing.T) {
	env1 := env.New([]string{"GO111MODULES=off"})

	env1.Set("GOROOT", "/go/1") // imaginary project 1
	env1.Set("GO111MODULES", "on")
	newStateManager(env1).SaveEnv()

	env2 := env.New(env1.Environ()) // next prompt

	env2.Set("GOROOT", "/go/2") // imaginary project 2
	env2.Unset("GO111MODULES")
	newStateManager(env2).SaveEnv()

	env3 := env.New(env2.Environ()) // next prompt

	newStateManager(env3).RestoreEnv()

	require.False(t, env3.Has("GOROOT"))
	require.Equal(t, "off", env3.Get("GO111MODULES"))
}
