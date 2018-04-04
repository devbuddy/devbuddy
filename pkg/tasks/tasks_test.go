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

func TestPython(t *testing.T) {
	task := ensureLoadTestTask(t, `python: 3.6.3`)

	require.Equal(t, task.(*Python).version, "3.6.3")
}
