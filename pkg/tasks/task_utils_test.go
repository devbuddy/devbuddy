package tasks

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func loadTestTask(t *testing.T, payload string) (*taskapi.Task, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(payload), &data)
	if err != nil {
		t.Fatalf("Failed to load a test fixture: %s", err)
	}

	task, err := taskapi.NewTaskFromDefinition(data)
	return task, err
}

func ensureLoadTestTask(t *testing.T, payload string) *taskapi.Task {
	task, err := loadTestTask(t, payload)
	require.NoError(t, err, "buildFromDefinition() failed")
	return task
}
