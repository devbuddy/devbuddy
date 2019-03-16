package autoenv

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/stretchr/testify/require"
)

func TestStateDeserialization(t *testing.T) {
	envs := [][]string{
		{},
		{"__BUD_AUTOENV_STATE="},
		{"__BUD_AUTOENV_STATE=f1=v1"},
		{"__BUD_AUTOENV_STATE=f1=v1:f2=v2"},

		{"__BUD_AUTOENV_STATE=1:project-1234:"},
		{"__BUD_AUTOENV_STATE=1:project-1234:f1=v1"},
		{"__BUD_AUTOENV_STATE=1:project-1234:f1=v1:f2=v2"},
	}
	features := []FeatureSet{
		NewFeatureSet(),
		NewFeatureSet(),
		NewFeatureSet().With(&FeatureInfo{"f1", "v1"}),
		NewFeatureSet().With(&FeatureInfo{"f1", "v1"}).With(&FeatureInfo{"f2", "v2"}),

		NewFeatureSet(),
		NewFeatureSet().With(&FeatureInfo{"f1", "v1"}),
		NewFeatureSet().With(&FeatureInfo{"f1", "v1"}).With(&FeatureInfo{"f2", "v2"}),
	}
	slugs := []string{
		"",
		"",
		"",
		"",

		"project-1234",
		"project-1234",
		"project-1234",
	}

	for idx := range envs {
		env := env.New(envs[idx])
		state := FeatureState{env: env}
		require.Equal(t, features[idx], state.GetActiveFeatures())
		require.Equal(t, slugs[idx], state.GetProjectSlug())
	}
}
func TestStateSerialization(t *testing.T) {
	env := env.New([]string{})
	state := FeatureState{env: env}

	state.SetProjectSlug("p-1")
	require.Equal(t, "1:p-1:", env.Get("__BUD_AUTOENV_STATE"))

	state.SetFeature(&FeatureInfo{"f1", "v1"})
	require.Equal(t, "1:p-1:f1=v1", env.Get("__BUD_AUTOENV_STATE"))

	state.SetProjectSlug("p-2")
	require.Equal(t, "1:p-2:f1=v1", env.Get("__BUD_AUTOENV_STATE"))

	state.UnsetFeature("f1")
	require.Equal(t, "1:p-2:", env.Get("__BUD_AUTOENV_STATE"))
}

func TestStateSetUnsetFeatures(t *testing.T) {
	env := env.New([]string{})

	newFeatureState := func() *FeatureState {
		return &FeatureState{env: env}
	}

	state := newFeatureState()
	state.SetFeature(&FeatureInfo{"rust", "v1"})

	state = newFeatureState()
	require.Equal(t,
		NewFeatureSet().With(&FeatureInfo{"rust", "v1"}),
		state.GetActiveFeatures())

	state = newFeatureState()
	state.SetFeature(&FeatureInfo{"elixir", "v1"})

	state = newFeatureState()
	require.Equal(t,
		NewFeatureSet().With(&FeatureInfo{"rust", "v1"}).With(&FeatureInfo{"elixir", "v1"}),
		state.GetActiveFeatures())

	state = newFeatureState()
	state.SetFeature(&FeatureInfo{"rust", "v2"})

	state = newFeatureState()
	require.Equal(t,
		NewFeatureSet().With(&FeatureInfo{"rust", "v2"}).With(&FeatureInfo{"elixir", "v1"}),
		state.GetActiveFeatures())

	state = newFeatureState()
	state.UnsetFeature("elixir")

	state = newFeatureState()
	require.Equal(t,
		NewFeatureSet().With(&FeatureInfo{"rust", "v2"}),
		state.GetActiveFeatures())

	state = newFeatureState()
	state.UnsetFeature("rust")

	state = newFeatureState()
	require.Equal(t,
		NewFeatureSet(),
		state.GetActiveFeatures())
}

func TestStateSetGetProjectSlug(t *testing.T) {
	env := env.New([]string{})

	newFeatureState := func() *FeatureState {
		return &FeatureState{env: env}
	}

	state := newFeatureState()
	state.SetProjectSlug("p-1")

	state = newFeatureState()
	require.Equal(t, "p-1", state.GetProjectSlug())

	state = newFeatureState()
	state.SetProjectSlug("p-123")

	state = newFeatureState()
	require.Equal(t, "p-123", state.GetProjectSlug())
}
