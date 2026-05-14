package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/stretchr/testify/require"
)

func withCwd(t *testing.T, path string) {
	t.Helper()

	previous, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(path))

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(previous))
	})
}

func TestLoadWithProject_ConfiguresExecutorFromProjectContext(t *testing.T) {
	projectDir := t.TempDir()
	resolvedProjectDir, err := filepath.EvalSymlinks(projectDir)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "dev.yml"), []byte(""), 0o600))
	withCwd(t, projectDir)

	ctx, err := LoadWithProject()
	require.NoError(t, err)

	require.NotNil(t, ctx.Project)
	require.Equal(t, resolvedProjectDir, ctx.Project.Path)
	require.NotNil(t, ctx.Executor)
	require.Equal(t, resolvedProjectDir, ctx.Executor.Cwd)
	require.Equal(t, "  ", ctx.Executor.OutputPrefix)
	require.Same(t, ctx.Env, ctx.Executor.Env)

	ctx.Env.Set("BUD_CONTEXT_TEST", "ok")
	result := ctx.Executor.Capture(executor.NewShell("echo $BUD_CONTEXT_TEST"))
	require.NoError(t, result.Error)
	require.Equal(t, "ok\n", result.Output)
}

func TestLoadWithProject_ReturnsProjectNotFound(t *testing.T) {
	withCwd(t, t.TempDir())

	ctx, err := LoadWithProject()
	require.ErrorIs(t, err, project.ErrProjectNotFound)
	require.Nil(t, ctx)
}
