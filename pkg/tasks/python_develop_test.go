package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPythonDevelop(t *testing.T) {
	task := ensureLoadTestTask(t, `
python_develop
`)

	require.Equal(t, "Task Python develop () has 1 actions", task.Describe())
}

func TestPythonDevelopWithExtras(t *testing.T) {
	task := ensureLoadTestTask(t, `
python_develop:
  extras: [dev, test]
`)

	require.Equal(t, "Task Python develop () has 1 actions", task.Describe())
}
