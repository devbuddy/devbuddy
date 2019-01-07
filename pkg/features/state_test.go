package features

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/stretchr/testify/require"
)

func TestStateSerialization(t *testing.T) {
	envs := [][]string{
		{},
		{"BUD_AUTO_ENV_FEATURES="},
		{"BUD_AUTO_ENV_FEATURES=f1=v1"},
		{"BUD_AUTO_ENV_FEATURES=f1=v1:f2=v2"},
	}
	features := []FeatureSet{
		NewFeatureSet(),
		NewFeatureSet(),
		NewFeatureSet().With(FeatureInfo{"f1", "v1"}),
		NewFeatureSet().With(FeatureInfo{"f1", "v1"}).With(FeatureInfo{"f2", "v2"}),
	}

	for idx := range envs {
		env := env.New(envs[idx])
		state := FeatureState{env}
		require.Equal(t, features[idx], state.GetActiveFeatures())
	}
}

func TestStateSetUnsetFeatures(t *testing.T) {
	env := env.New([]string{})

	state := FeatureState{env}
	state.SetFeature(FeatureInfo{"rust", "v1"})

	state = FeatureState{env}
	require.Equal(t,
		NewFeatureSet().With(FeatureInfo{"rust", "v1"}),
		state.GetActiveFeatures())

	state = FeatureState{env}
	state.SetFeature(FeatureInfo{"elixir", "v1"})

	state = FeatureState{env}
	require.Equal(t,
		NewFeatureSet().With(FeatureInfo{"rust", "v1"}).With(FeatureInfo{"elixir", "v1"}),
		state.GetActiveFeatures())

	state = FeatureState{env}
	state.SetFeature(FeatureInfo{"rust", "v2"})

	state = FeatureState{env}
	require.Equal(t,
		NewFeatureSet().With(FeatureInfo{"rust", "v2"}).With(FeatureInfo{"elixir", "v1"}),
		state.GetActiveFeatures())

	state = FeatureState{env}
	state.UnsetFeature("elixir")

	state = FeatureState{env}
	require.Equal(t,
		NewFeatureSet().With(FeatureInfo{"rust", "v2"}),
		state.GetActiveFeatures())

	state = FeatureState{env}
	state.UnsetFeature("rust")

	state = FeatureState{env}
	require.Equal(t,
		NewFeatureSet(),
		state.GetActiveFeatures())
}
