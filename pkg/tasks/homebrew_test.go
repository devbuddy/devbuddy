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

func TestHomebrewEmpty(t *testing.T) {
	_, err := loadTestTask(t, `homebrew: []`)
	require.EqualError(t, err, `task "homebrew": no homebrew formulas specified`)
}

func TestHomebrewInvalidElementType(t *testing.T) {
	_, err := loadTestTask(t, `
homebrew:
  - git
  - false
`)
	require.EqualError(t, err, `task "homebrew": expecting a list of strings, found an invalid element: type bool (false)`)
}
