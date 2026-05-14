package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPythonDevelop(t *testing.T) {
	task := ensureLoadTestTask(t, `
python_develop
`)

	require.Equal(t, "Task Python develop () required_task=python actions=1", task.Describe())
	require.Equal(t, "python", task.RequiredTask)
	require.Equal(t, 1, len(task.Actions))
	requireNoFeature(t, task)
}

func TestPythonDevelopWithExtras(t *testing.T) {
	task := ensureLoadTestTask(t, `
python_develop:
  extras: [dev, test]
`)

	require.Equal(t, "Task Python develop () required_task=python actions=1", task.Describe())
	require.Equal(t, "python", task.RequiredTask)
	require.Equal(t, 1, len(task.Actions))
	requireNoFeature(t, task)
}
