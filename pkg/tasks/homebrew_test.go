package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHomebrew(t *testing.T) {
	task := ensureLoadTestTask(t, `
homebrew:
  - file1
  - file2
`)
	require.Equal(t, "Task Homebrew (file1, file2) actions=2", task.Describe())
}
