package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPythonDevelop(t *testing.T) {
	task := ensureLoadTestTask(t, `
python_develop
`)

	require.Equal(t, taskDefinitions["python_develop"], task.taskDefinition)
	require.Equal(t, 1, len(task.actions))
	require.Equal(t, "install python package in develop mode", task.actions[0].description())
}

func TestPythonDevelopWithExtras(t *testing.T) {
	task := ensureLoadTestTask(t, `
python_develop:
  extras: [dev, test]
`)

	require.Equal(t, taskDefinitions["python_develop"], task.taskDefinition)
	require.Equal(t, 1, len(task.actions))
	require.Equal(t, "install python package in develop mode", task.actions[0].description())

}
