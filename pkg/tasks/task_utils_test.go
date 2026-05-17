package tasks

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	yaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

func loadTestTask(t *testing.T, payload string) (*api.Task, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(payload), &data)
	require.NoError(t, err, "Failed to load a yaml task fixture: %s")
	return api.NewTaskFromPayload(data)
}

func ensureLoadTestTask(t *testing.T, payload string) *api.Task {
	task, err := loadTestTask(t, payload)
	require.NoError(t, err, "NewTaskFromPayload() failed")
	return task
}

// taskFeature returns the first feature found across all actions, or nil.
func taskFeature(task *api.Task) *autoenv.FeatureInfo {
	for _, action := range task.Actions {
		if f := action.Feature(); f != nil {
			return f
		}
	}
	return nil
}

// requireFeature asserts that the task has a feature with the given name and param.
func requireFeature(t *testing.T, task *api.Task, expectedName, expectedParam string) {
	t.Helper()
	f := taskFeature(task)
	require.NotNil(t, f, "expected a feature but none found")
	require.Equal(t, expectedName, f.Name, "feature name")
	require.Equal(t, expectedParam, f.Param, "feature param")
}

// requireNoFeature asserts that no action on the task declares a feature.
func requireNoFeature(t *testing.T, task *api.Task) {
	t.Helper()
	for _, action := range task.Actions {
		require.Nil(t, action.Feature(), "expected no feature but found one")
	}
}
