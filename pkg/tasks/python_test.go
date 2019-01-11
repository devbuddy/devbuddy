package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPythonOk(t *testing.T) {
	task := ensureLoadTestTask(t, `python: 3.6.3`)

	require.Equal(t, "Task Python (3.6.3) has feature python:3.6.3 and has 4 actions", task.Describe())
}

func TestPythonInvalid(t *testing.T) {
	_, err := loadTestTask(t, `python: 3`)

	require.Error(t, err, "buildFromDefinition() should have failed")
}
