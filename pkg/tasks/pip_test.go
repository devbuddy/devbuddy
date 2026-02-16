package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPip(t *testing.T) {
	task := ensureLoadTestTask(t, `
pip:
  - file1
  - file2
`)

	require.Equal(t, "Task Pip (file1, file2) required_task=python actions=2", task.Describe())
	require.Equal(t, "file1, file2", task.Info)
	require.Equal(t, "python", task.RequiredTask)
	require.Equal(t, 2, len(task.Actions))
}

func TestPipEmpty(t *testing.T) {
	_, err := loadTestTask(t, `pip: []`)
	require.EqualError(t, err, `task "pip": no pip files specified`)
}

func TestPipInvalidElementType(t *testing.T) {
	_, err := loadTestTask(t, `
pip:
  - requirements.txt
  - 123
`)
	require.EqualError(t, err, `task "pip": expecting a list of strings, found an invalid element: type int (123)`)
}
