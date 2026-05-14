package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRubyOk(t *testing.T) {
	task := ensureLoadTestTask(t, `ruby: 3.3.0`)

	require.Equal(t, "Task Ruby (3.3.0) feature=ruby:3.3.0 actions=3", task.Describe())
	require.Equal(t, "3.3.0", task.Info)
	require.Equal(t, 3, len(task.Actions))
	requireFeature(t, task, "ruby", "3.3.0")
}

func TestRubyInvalid(t *testing.T) {
	_, err := loadTestTask(t, `ruby: 3`)

	require.Error(t, err, "buildFromDefinition() should have failed")
}
