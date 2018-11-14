package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApt(t *testing.T) {
	task := ensureLoadTestTask(t, `
apt:
  - curl
  - git
`)

	require.Equal(t, "curl, git", task.header)
	require.Equal(t, 1, len(task.actions))
}

func TestAptEmpty(t *testing.T) {
	_, err := loadTestTask(t, `
apt: []
`)
	require.Error(t, err)
}
