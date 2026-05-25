package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUV(t *testing.T) {
	task := ensureLoadTestTask(t, `uv`)

	require.Equal(t, "Task uv (sync) required_task=python actions=2", task.Describe())
	require.Equal(t, "sync", task.Info)
	require.Equal(t, "python", task.RequiredTask)
	require.Equal(t, 2, len(task.Actions))
	require.Equal(t, "install uv", task.Actions[0].Description())
	require.Equal(t, "sync uv project", task.Actions[1].Description())
}

func TestUVSyncArgsDefaultToActiveInexactEnvironment(t *testing.T) {
	args := buildUVSyncArgs(uvConfig{})

	require.Equal(t, []string{"sync", "--active", "--inexact"}, args)
}

func TestUVSyncArgs(t *testing.T) {
	args := buildUVSyncArgs(uvConfig{
		Groups:    []string{"dev", "docs"},
		Extras:    []string{"postgres"},
		AllGroups: true,
		AllExtras: true,
		Exact:     true,
		Frozen:    true,
	})

	require.Equal(t, []string{
		"sync",
		"--active",
		"--group", "dev",
		"--group", "docs",
		"--extra", "postgres",
		"--all-groups",
		"--all-extras",
		"--frozen",
	}, args)
}
