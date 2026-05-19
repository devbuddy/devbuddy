package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPythonOk(t *testing.T) {
	task := ensureLoadTestTask(t, `python: 3.6.3`)

	require.Equal(t, "Task Python (3.6.3) feature=python:3.6.3 actions=3", task.Describe())
	require.Equal(t, "3.6.3", task.Info)
	require.Equal(t, 3, len(task.Actions))
	require.Equal(t, "install PyEnv", task.Actions[0].Description())
	require.Equal(t, "install Python version with PyEnv", task.Actions[1].Description())
	require.Equal(t, "create virtualenv", task.Actions[2].Description())
	requireFeature(t, task, "python", "3.6.3")
}

func TestPythonInvalid(t *testing.T) {
	_, err := loadTestTask(t, `python: 3`)

	require.Error(t, err, "buildFromDefinition() should have failed")
}
