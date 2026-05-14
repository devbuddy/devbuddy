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

	require.Equal(t, "Task Apt (curl, git) actions=1", task.Describe())
	require.Equal(t, "curl, git", task.Info)
	require.Equal(t, 1, len(task.Actions))
	requireNoFeature(t, task)
}

func TestAptEmpty(t *testing.T) {
	_, err := loadTestTask(t, `
apt: []
`)
	require.Error(t, err)
}
