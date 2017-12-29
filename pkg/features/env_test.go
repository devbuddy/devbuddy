package features

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetActiveFeatures(t *testing.T) {
	envs := [][]string{
		[]string{},
		[]string{"DEV_AUTO_ENV_FEATURES="},
		[]string{"DEV_AUTO_ENV_FEATURES=f1=v1"},
		[]string{"DEV_AUTO_ENV_FEATURES=f1=v1:f2=v2"},
	}
	features := []map[string]string{
		map[string]string{},
		map[string]string{},
		map[string]string{"f1": "v1"},
		map[string]string{"f1": "v1", "f2": "v2"},
	}

	for idx := range envs {
		env := NewEnv(envs[idx])
		require.Equal(t, features[idx], env.GetActiveFeatures())
	}
}

func TestSetFeatures(t *testing.T) {
	features := []map[string]string{
		map[string]string{},
		map[string]string{},
		map[string]string{"f1": "v1"},
		map[string]string{"f1": "v1", "f2": "v2"},
	}

	for idx := range features {
		env := NewEnv([]string{})
		env.SetActiveFeatures(features[idx])
		require.Equal(t, features[idx], env.GetActiveFeatures())
	}
}

func TestChanged(t *testing.T) {
	env := NewEnv([]string{})
	require.Equal(t, []EnvVarChange{}, env.Changed())

	env.Set("K2", "1")
	require.Equal(t, []EnvVarChange{EnvVarChange{"K2", "1", false}}, env.Changed())

	env.Set("K2", "2")
	require.Equal(t, []EnvVarChange{EnvVarChange{"K2", "2", false}}, env.Changed())

	env = NewEnv([]string{"K1=1"})
	env.Unset("K1")
	require.Equal(t, []EnvVarChange{EnvVarChange{"K1", "", true}}, env.Changed())
}
