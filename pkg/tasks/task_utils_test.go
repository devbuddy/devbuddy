package tasks

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
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
