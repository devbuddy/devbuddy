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

	require.Equal(t, task.(*Pip).files, []string{"file1", "file2"})
}
