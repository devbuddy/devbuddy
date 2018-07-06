package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func loadTestTask(t *testing.T, payload string) (*Task, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(payload), &data)
	if err != nil {
		t.Fatalf("Failed to load a test fixture: %s", err)
	}

	task, err := buildFromDefinition(data)
	return task, err
}

func ensureLoadTestTask(t *testing.T, payload string) *Task {
	task, err := loadTestTask(t, payload)
	require.NoError(t, err, "buildFromDefinition() failed")
	return task
}
