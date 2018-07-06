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

	require.Equal(t, len(task.actions), 2)
	require.Equal(t, task.header, "file1, file2")
}
