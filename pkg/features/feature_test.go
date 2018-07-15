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
	require.ElementsMatch(t, []string{"python", "golang"}, names)

	for _, name := range names {
		require.NotNil(t, definitions.Get(name))
	}
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

	devUpNeeded, err = Refresh("test-activation", "testparam", nil, nil, nil)
	require.False(t, devUpNeeded)
	require.NoError(t, err)

	Deactivate("test-activation", "testparam", nil, nil)
}

func TestRefresh(t *testing.T) {
	refreshed := false

	d := definitions.Register("test-refresh")
	d.Refresh = func(param string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
		require.Equal(t, "testparam", param)
		refreshed = true
		return true, nil
	}

	devUpNeeded, err := Activate("test-refresh", "testparam", nil, nil, nil)
	require.False(t, devUpNeeded)
	require.NoError(t, err)

	devUpNeeded, err = Refresh("test-refresh", "testparam", nil, nil, nil)
	require.True(t, devUpNeeded)
	require.NoError(t, err)
	require.True(t, refreshed)

	Deactivate("test-refresh", "testparam", nil, nil)
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

	devUpNeeded, err = Refresh("test-deactivate", "testparam", nil, nil, nil)
	require.False(t, devUpNeeded)
	require.NoError(t, err)

	Deactivate("test-deactivate", "testparam", nil, nil)
	require.True(t, deactivated)
}
