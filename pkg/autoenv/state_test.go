package autoenv

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/stretchr/testify/require"
)

func TestStateDeserialization(t *testing.T) {
	envs := [][]string{
		{},
		{"BUD_AUTO_ENV_FEATURES="},
		{"BUD_AUTO_ENV_FEATURES=f1=v1"},
		{"BUD_AUTO_ENV_FEATURES=f1=v1:f2=v2"},

		{"BUD_AUTO_ENV_FEATURES=1:project-1234:"},
		{"BUD_AUTO_ENV_FEATURES=1:project-1234:f1=v1"},
		{"BUD_AUTO_ENV_FEATURES=1:project-1234:f1=v1:f2=v2"},
	}
	features := []FeatureSet{
		NewFeatureSet(),
		NewFeatureSet(),
		NewFeatureSet().With(NewFeatureInfo("f1", "v1")),
		NewFeatureSet().With(NewFeatureInfo("f1", "v1")).With(NewFeatureInfo("f2", "v2")),

		NewFeatureSet(),
		NewFeatureSet().With(NewFeatureInfo("f1", "v1")),
		NewFeatureSet().With(NewFeatureInfo("f1", "v1")).With(NewFeatureInfo("f2", "v2")),
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
		state := FeatureState{env}
		require.Equal(t, features[idx], state.GetActiveFeatures())
		require.Equal(t, slugs[idx], state.GetProjectSlug())
	}
}
func TestStateSerialization(t *testing.T) {
	env := env.New([]string{})
	state := FeatureState{env}

	state.SetProjectSlug("p-1")
	require.Equal(t, "1:p-1:", env.Get("BUD_AUTO_ENV_FEATURES"))

	state.SetFeature(NewFeatureInfo("f1", "v1"))
	require.Equal(t, "1:p-1:f1=v1", env.Get("BUD_AUTO_ENV_FEATURES"))

	state.SetProjectSlug("p-2")
	require.Equal(t, "1:p-2:f1=v1", env.Get("BUD_AUTO_ENV_FEATURES"))

	state.UnsetFeature("f1")
	require.Equal(t, "1:p-2:", env.Get("BUD_AUTO_ENV_FEATURES"))
}

func TestStateSetUnsetFeatures(t *testing.T) {
	env := env.New([]string{})

	state := FeatureState{env}
	state.SetFeature(NewFeatureInfo("rust", "v1"))

	state = FeatureState{env}
	require.Equal(t, "rust:v1", state.GetActiveFeatures().String())

	state = FeatureState{env}
	state.SetFeature(NewFeatureInfo("elixir", "v1"))

	state = FeatureState{env}
	require.Equal(t, "rust:v1 elixir:v1", state.GetActiveFeatures().String())

	state = FeatureState{env}
	state.SetFeature(NewFeatureInfo("rust", "v2"))

	state = FeatureState{env}
	require.Equal(t, "elixir:v1 rust:v2", state.GetActiveFeatures().String())

	state = FeatureState{env}
	state.UnsetFeature("elixir")

	state = FeatureState{env}
	require.Equal(t, "rust:v2", state.GetActiveFeatures().String())

	state = FeatureState{env}
	state.UnsetFeature("rust")

	state = FeatureState{env}
	require.Equal(t, "", state.GetActiveFeatures().String())
}

func TestStateSetGetProjectSlug(t *testing.T) {
	env := env.New([]string{})

	state := FeatureState{env}
	state.SetProjectSlug("p-1")

	state = FeatureState{env}
	require.Equal(t, "p-1", state.GetProjectSlug())

	state = FeatureState{env}
	state.SetProjectSlug("p-123")

	state = FeatureState{env}
	require.Equal(t, "p-123", state.GetProjectSlug())
}
