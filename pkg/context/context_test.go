package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/ui"
	"github.com/stretchr/testify/require"
)

type taskCommandRunnerSpy struct {
	runCmd *executor.Command
}

func (s *taskCommandRunnerSpy) Run(cmd *executor.Command) *executor.Result {
	s.runCmd = cmd
	return &executor.Result{}
}

func (s *taskCommandRunnerSpy) Capture(cmd *executor.Command) *executor.Result {
	return &executor.Result{}
}

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

func TestRunTaskCommand_DisplaysAndRunsSameCommand(t *testing.T) {
	_, testingUI := ui.NewTesting()
	runner := &taskCommandRunnerSpy{}
	ctx := &Context{
		UI:       testingUI,
		Executor: &executor.Executor{Runner: runner},
	}
	cmd := executor.New("pip", "install", "-r", "requirements.txt").AddOutputFilter("already satisfied")

	result := ctx.RunTaskCommand(cmd)

	require.NoError(t, result.Error)
	require.Same(t, cmd, runner.runCmd)

	events := testingUI.Events()
	require.Len(t, events, 1)
	require.Equal(t, ui.KindTaskCommand, events[0].Kind)
	require.Equal(t, "pip", events[0].Text)
	require.Equal(t, []ui.Field{ui.F("args", "install -r requirements.txt")}, events[0].Fields)
}
