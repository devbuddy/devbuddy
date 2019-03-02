package features

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"

	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	reg := NewRegister()

	activate := func(param string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
		return false, nil
	}
	deactivate1 := func(param string, cfg *config.Config, env *env.Env) {}
	deactivate2 := func(param string, cfg *config.Config, env *env.Env) {}

	reg.Register("env1", activate, deactivate1)
	reg.Register("env2", activate, deactivate2)

	require.ElementsMatch(t, reg.Names(), []string{"env1", "env2"})

	env, err := reg.Get("env1")
	require.NoError(t, err)
	require.Equal(t, env.Name, "env1")
	require.NotNil(t, env.Activate)
	require.NotNil(t, env.Deactivate)

	env, err = reg.Get("env2")
	require.NoError(t, err)
	require.Equal(t, env.Name, "env2")
	require.NotNil(t, env.Activate)
	require.NotNil(t, env.Deactivate)
}

func TestRegisterNotFound(t *testing.T) {
	reg := NewRegister()

	_, err := reg.Get("nope")
	require.Error(t, err)
}
