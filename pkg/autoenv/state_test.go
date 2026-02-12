package autoenv

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/termui"

	"github.com/stretchr/testify/require"
)

func newStateManager(env *env.Env) *StateManager {
	_, ui := termui.NewTesting(false)
	return &StateManager{env: env, UI: ui}
}

// mustGetActiveFeatures calls GetActiveFeatures and fails the test on error.
func mustGetActiveFeatures(t *testing.T, s *StateManager) FeatureSet {
	t.Helper()
	fs, err := s.GetActiveFeatures()
	require.NoError(t, err)
	return fs
}

// mustGetProjectSlug calls GetProjectSlug and fails the test on error.
func mustGetProjectSlug(t *testing.T, s *StateManager) string {
	t.Helper()
	slug, err := s.GetProjectSlug()
	require.NoError(t, err)
	return slug
}

func TestStateSerialization(t *testing.T) {
	env := env.New([]string{})
	state := newStateManager(env)

	require.NoError(t, state.SetProjectSlug("p-1"))
	require.Equal(t, `{"project":"p-1","features":null,"saved_env":{}}`,
		env.Get("__BUD_AUTOENV"))

	require.NoError(t, state.SetFeature(&FeatureInfo{"f1", "v1"}))
	require.Equal(t, `{"project":"p-1","features":[{"name":"f1","param":"v1"}],"saved_env":{}}`,
		env.Get("__BUD_AUTOENV"))

	require.NoError(t, state.SetProjectSlug("p-2"))
	require.Equal(t, `{"project":"p-2","features":[{"name":"f1","param":"v1"}],"saved_env":{}}`,
		env.Get("__BUD_AUTOENV"))

	require.NoError(t, state.UnsetFeature("f1"))
	require.Equal(t, `{"project":"p-2","features":[],"saved_env":{}}`,
		env.Get("__BUD_AUTOENV"))
}

func TestStateSetUnsetFeatures(t *testing.T) {
	env := env.New([]string{})

	require.NoError(t, newStateManager(env).SetFeature(NewFeatureInfo("rust", "v1")))
	require.Equal(t, "rust:v1", mustGetActiveFeatures(t, newStateManager(env)).String())

	require.NoError(t, newStateManager(env).SetFeature(NewFeatureInfo("elixir", "v1")))
	require.Equal(t, "rust:v1 elixir:v1", mustGetActiveFeatures(t, newStateManager(env)).String())

	require.NoError(t, newStateManager(env).SetFeature(NewFeatureInfo("rust", "v2")))
	require.Equal(t, "elixir:v1 rust:v2", mustGetActiveFeatures(t, newStateManager(env)).String())

	require.NoError(t, newStateManager(env).UnsetFeature("elixir"))
	require.Equal(t, "rust:v2", mustGetActiveFeatures(t, newStateManager(env)).String())

	require.NoError(t, newStateManager(env).UnsetFeature("rust"))
	require.Equal(t, "", mustGetActiveFeatures(t, newStateManager(env)).String())
}

func TestStateSetGetProjectSlug(t *testing.T) {
	env := env.New([]string{})

	require.NoError(t, newStateManager(env).SetProjectSlug("p-1"))
	require.Equal(t, "p-1", mustGetProjectSlug(t, newStateManager(env)))

	require.NoError(t, newStateManager(env).SetProjectSlug("p-123"))
	require.Equal(t, "p-123", mustGetProjectSlug(t, newStateManager(env)))
}

func TestStateSavedEnvDirect(t *testing.T) {
	env1 := env.New([]string{"GO111MODULES=off"})

	env1.Set("GOROOT", "/go/1") // imaginary project
	env1.Set("GO111MODULES", "on")
	require.NoError(t, newStateManager(env1).SaveEnv())

	env2 := env.New(env1.Environ()) // next prompt

	require.NoError(t, newStateManager(env2).RestoreEnv())

	require.False(t, env2.Has("GOROOT"))
	require.Equal(t, "off", env2.Get("GO111MODULES"))
}

func TestStateSavedEnvMultipleChanges(t *testing.T) {
	env1 := env.New([]string{"GO111MODULES=off"})

	env1.Set("GOROOT", "/go/1") // imaginary project 1
	env1.Set("GO111MODULES", "on")
	require.NoError(t, newStateManager(env1).SaveEnv())

	env2 := env.New(env1.Environ()) // next prompt

	env2.Set("GOROOT", "/go/2") // imaginary project 2
	env2.Unset("GO111MODULES")
	require.NoError(t, newStateManager(env2).SaveEnv())

	env3 := env.New(env2.Environ()) // next prompt

	require.NoError(t, newStateManager(env3).RestoreEnv())

	require.False(t, env3.Has("GOROOT"))
	require.Equal(t, "off", env3.Get("GO111MODULES"))
}

func TestStateSavedEnvForget(t *testing.T) {
	env1 := env.New([]string{})

	env1.Set("GOROOT", "/go/1")
	require.NoError(t, newStateManager(env1).SaveEnv())

	env2 := env.New(env1.Environ())

	stateManager := newStateManager(env2)
	require.NoError(t, stateManager.ForgetEnv())
	require.NoError(t, stateManager.RestoreEnv()) // the restore would have unset GOROOT

	require.True(t, env2.Has("GOROOT"))
}

func TestStateInvalidJSON(t *testing.T) {
	e := env.New([]string{"__BUD_AUTOENV=not-valid-json"})
	state := newStateManager(e)

	_, err := state.GetActiveFeatures()
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read the state")
}
