package features

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	names := definitions.Names()
	require.Contains(t, names, "python")
	require.Contains(t, names, "golang")
}

func TestActivation(t *testing.T) {
	activated := false

	d := definitions.Register("test-activation")
	d.Activate = func(param string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
		require.Equal(t, "testparam", param)
		activated = true
		return true, nil
	}

	devUpNeeded, err := Activate("test-activation", "testparam", nil, nil, nil)
	require.True(t, devUpNeeded)
	require.NoError(t, err)
	require.True(t, activated)

	Deactivate("test-activation", "testparam", nil, nil)
}

func TestDeactivate(t *testing.T) {
	deactivated := false

	d := definitions.Register("test-deactivate")
	d.Deactivate = func(param string, cfg *config.Config, env *env.Env) {
		require.Equal(t, "testparam", param)
		deactivated = true
	}

	devUpNeeded, err := Activate("test-deactivate", "testparam", nil, nil, nil)
	require.False(t, devUpNeeded)
	require.NoError(t, err)

	Deactivate("test-deactivate", "testparam", nil, nil)
	require.True(t, deactivated)
}
