package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPythonOk(t *testing.T) {
	task := ensureLoadTestTask(t, `python: 3.6.3`)

	require.Equal(t, "Task Python (3.6.3) feature=python:3.6.3 actions=4", task.Describe())
	require.Equal(t, "3.6.3", task.Info)
	require.Equal(t, 4, len(task.Actions))
	requireFeature(t, task, "python", "3.6.3")
}

func TestPythonInvalid(t *testing.T) {
	_, err := loadTestTask(t, `python: 3`)

	require.Error(t, err, "buildFromDefinition() should have failed")
}
