package autoenv

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/stretchr/testify/require"
)

func TestStateSerialization(t *testing.T) {
	env := env.New([]string{})
	state := &StateManager{env: env}

	state.SetProjectSlug("p-1")
	require.Equal(t, `{"project":"p-1","features":null,"saved_state":{}}`,
		env.Get("__BUD_AUTOENV"))

	state.SetFeature(&FeatureInfo{"f1", "v1"})
	require.Equal(t, `{"project":"p-1","features":[{"name":"f1","param":"v1"}],"saved_state":{}}`,
		env.Get("__BUD_AUTOENV"))

	state.SetProjectSlug("p-2")
	require.Equal(t, `{"project":"p-2","features":[{"name":"f1","param":"v1"}],"saved_state":{}}`,
		env.Get("__BUD_AUTOENV"))

	state.UnsetFeature("f1")
	require.Equal(t, `{"project":"p-2","features":[],"saved_state":{}}`,
		env.Get("__BUD_AUTOENV"))
}

func TestStateSetUnsetFeatures(t *testing.T) {
	env := env.New([]string{})

	newStateManager := func() *StateManager {
		return &StateManager{env: env}
	}

	newStateManager().SetFeature(NewFeatureInfo("rust", "v1"))
	require.Equal(t, "rust:v1", newStateManager().GetActiveFeatures().String())

	newStateManager().SetFeature(NewFeatureInfo("elixir", "v1"))
	require.Equal(t, "rust:v1 elixir:v1", newStateManager().GetActiveFeatures().String())

	newStateManager().SetFeature(NewFeatureInfo("rust", "v2"))
	require.Equal(t, "elixir:v1 rust:v2", newStateManager().GetActiveFeatures().String())

	newStateManager().UnsetFeature("elixir")
	require.Equal(t, "rust:v2", newStateManager().GetActiveFeatures().String())

	newStateManager().UnsetFeature("rust")
	require.Equal(t, "", newStateManager().GetActiveFeatures().String())
}

func TestStateSetGetProjectSlug(t *testing.T) {
	env := env.New([]string{})

	newStateManager := func() *StateManager {
		return &StateManager{env: env}
	}

	newStateManager().SetProjectSlug("p-1")
	require.Equal(t, "p-1", newStateManager().GetProjectSlug())

	newStateManager().SetProjectSlug("p-123")
	require.Equal(t, "p-123", newStateManager().GetProjectSlug())
}
