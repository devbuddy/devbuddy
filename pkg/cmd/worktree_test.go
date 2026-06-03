package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/ui"
	"github.com/devbuddy/devbuddy/pkg/worktree"
	"github.com/stretchr/testify/require"
)

func TestFormatWorktreeRowsAlignsColumns(t *testing.T) {
	rows := []worktreeRow{
		{
			Worktree: worktree.Worktree{Path: "/src/github.com/acme/api", Branch: "main"},
			Branch:   "main",
			Head:     "1111111",
			State:    "clean",
			Modified: "2026-05-15",
		},
		{
			Worktree: worktree.Worktree{Path: "/src/github.com/acme/api--agent-1", Branch: "long-feature-branch"},
			Branch:   "long-feature-branch",
			Head:     "2222222",
			State:    "dirty",
			Modified: "2026-05-15",
		},
	}

	lines := formatWorktreeRows(rows, false)

	require.Len(t, lines, 2)
	require.Equal(t, strings.Index(lines[0], "/src/"), strings.Index(lines[1], "/src/"))
	require.Contains(t, lines[0], "main")
	require.Contains(t, lines[1], "long-feature-branch")
}

func TestSelectWorktreeUsesPromptOptions(t *testing.T) {
	dir := t.TempDir()
	firstPath := filepath.Join(dir, "api")
	secondPath := filepath.Join(dir, "api--feature-a")
	require.NoError(t, os.Mkdir(firstPath, 0755))
	require.NoError(t, os.Mkdir(secondPath, 0755))

	prompts := &ui.FakePrompts{SelectValue: secondPath}
	exec := &executor.Executor{Runner: worktreeRunner{}}

	got, err := selectWorktree(prompts, exec, []worktree.Worktree{
		{Path: firstPath, Branch: "main", Head: "111111111111"},
		{Path: secondPath, Branch: "feature-a", Head: "222222222222"},
	})

	require.NoError(t, err)
	require.Equal(t, secondPath, got.Path)
	require.Len(t, prompts.SelectRequests, 1)
	require.Equal(t, "Select worktree", prompts.SelectRequests[0].Label)
	require.Equal(t, []ui.SelectOption{
		{Value: firstPath, Label: prompts.SelectRequests[0].Options[0].Label},
		{Value: secondPath, Label: prompts.SelectRequests[0].Options[1].Label},
	}, prompts.SelectRequests[0].Options)
	require.Contains(t, prompts.SelectRequests[0].Options[0].Label, "main")
	require.Contains(t, prompts.SelectRequests[0].Options[1].Label, "feature-a")
}

func TestSelectWorktreeReturnsPromptCancellation(t *testing.T) {
	dir := t.TempDir()
	worktreePath := filepath.Join(dir, "api")
	require.NoError(t, os.Mkdir(worktreePath, 0755))

	prompts := &ui.FakePrompts{SelectErr: ui.ErrPromptCancelled}
	exec := &executor.Executor{Runner: worktreeRunner{}}

	_, err := selectWorktree(prompts, exec, []worktree.Worktree{
		{Path: worktreePath, Branch: "main", Head: "111111111111"},
	})

	require.ErrorIs(t, err, ui.ErrPromptCancelled)
}

func TestInactiveWorktreesSkipsMainAndRecentWorktrees(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "api")
	oldPath := filepath.Join(dir, "api--old")
	recentPath := filepath.Join(dir, "api--recent")
	require.NoError(t, os.Mkdir(mainPath, 0755))
	require.NoError(t, os.Mkdir(oldPath, 0755))
	require.NoError(t, os.Mkdir(recentPath, 0755))

	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	require.NoError(t, os.Chtimes(mainPath, now.Add(-14*24*time.Hour), now.Add(-14*24*time.Hour)))
	require.NoError(t, os.Chtimes(oldPath, now.Add(-8*24*time.Hour), now.Add(-8*24*time.Hour)))
	require.NoError(t, os.Chtimes(recentPath, now.Add(-2*24*time.Hour), now.Add(-2*24*time.Hour)))

	got := inactiveWorktrees([]worktree.Worktree{
		{Path: mainPath, Branch: "main"},
		{Path: oldPath, Branch: "old"},
		{Path: recentPath, Branch: "recent"},
	}, now, 7*24*time.Hour)

	require.Equal(t, []worktree.Worktree{{Path: oldPath, Branch: "old"}}, got)
}

type worktreeRunner struct{}

func (worktreeRunner) Run(*executor.Command) *executor.Result {
	return &executor.Result{}
}

func (worktreeRunner) Capture(*executor.Command) *executor.Result {
	return &executor.Result{}
}

func TestWorktreeRemoveRun(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "api")
	otherPath := filepath.Join(dir, "api--feat")
	require.NoError(t, os.Mkdir(mainPath, 0755))
	require.NoError(t, os.Mkdir(otherPath, 0755))

	worktrees := []worktree.Worktree{
		{Path: mainPath, Branch: "main"},
		{Path: otherPath, Branch: "feat"},
	}

	newCtx := func(t *testing.T, projectPath string, prompts ui.Prompts) (*context.Context, string) {
		t.Helper()
		_, buf := ui.NewBufferedTesting(false)
		buf.SetPrompts(prompts)
		ctx := &context.Context{
			Project:  &project.Project{Path: projectPath},
			Executor: &executor.Executor{Runner: worktreeRunner{}},
			UI:       buf,
		}
		finalizer := filepath.Join(t.TempDir(), "finalizer")
		require.NoError(t, os.WriteFile(finalizer, []byte(""), 0644))
		t.Setenv("BUD_FINALIZER_FILE", finalizer)
		return ctx, finalizer
	}

	listFn := func(*executor.Executor, string) ([]worktree.Worktree, error) {
		return worktrees, nil
	}

	t.Run("refuses to remove the main worktree", func(t *testing.T) {
		prompts := &ui.FakePrompts{ConfirmValue: true}
		ctx, _ := newCtx(t, mainPath, prompts)

		err := runWorktreeRemove(ctx, []string{"main"}, listFn)

		require.Error(t, err)
		require.Contains(t, err.Error(), "refusing to remove the main worktree")
		require.Empty(t, prompts.ConfirmRequests, "confirmation should not be asked for the main worktree")
	})

	t.Run("defaults to the current worktree", func(t *testing.T) {
		prompts := &ui.FakePrompts{ConfirmValue: true}
		ctx, _ := newCtx(t, otherPath, prompts)

		err := runWorktreeRemove(ctx, nil, listFn)

		require.NoError(t, err)
		require.Len(t, prompts.ConfirmRequests, 1)
		require.Contains(t, prompts.ConfirmRequests[0].Label, "feat")
		require.Contains(t, prompts.ConfirmRequests[0].Label, otherPath)
	})

	t.Run("errors when defaulting from the main worktree", func(t *testing.T) {
		prompts := &ui.FakePrompts{ConfirmValue: true}
		ctx, _ := newCtx(t, mainPath, prompts)

		err := runWorktreeRemove(ctx, nil, listFn)

		require.Error(t, err)
		require.Contains(t, err.Error(), "refusing to remove the main worktree")
		require.Empty(t, prompts.ConfirmRequests)
	})

	t.Run("errors when current directory is not a tracked worktree", func(t *testing.T) {
		unknownPath := filepath.Join(dir, "api--unknown")
		require.NoError(t, os.Mkdir(unknownPath, 0755))
		prompts := &ui.FakePrompts{ConfirmValue: true}
		ctx, _ := newCtx(t, unknownPath, prompts)

		err := runWorktreeRemove(ctx, nil, listFn)

		require.Error(t, err)
		require.Contains(t, err.Error(), "not a git worktree")
		require.Empty(t, prompts.ConfirmRequests)
	})

	t.Run("removes the worktree matched by query after confirmation", func(t *testing.T) {
		prompts := &ui.FakePrompts{ConfirmValue: true}
		ctx, _ := newCtx(t, mainPath, prompts)

		err := runWorktreeRemove(ctx, []string{"feat"}, listFn)

		require.NoError(t, err)
		require.Len(t, prompts.ConfirmRequests, 1)
		require.Contains(t, prompts.ConfirmRequests[0].Label, "feat")
	})

	t.Run("aborts when the user declines confirmation", func(t *testing.T) {
		prompts := &ui.FakePrompts{ConfirmValue: false}
		ctx, finalizer := newCtx(t, mainPath, prompts)

		err := runWorktreeRemove(ctx, []string{"feat"}, listFn)

		require.NoError(t, err)
		require.Len(t, prompts.ConfirmRequests, 1)
		content, readErr := os.ReadFile(finalizer)
		require.NoError(t, readErr)
		require.Empty(t, content, "no finalizer should be written when the removal is aborted")
	})

	t.Run("aborts when the confirmation prompt is cancelled", func(t *testing.T) {
		prompts := &ui.FakePrompts{ConfirmErr: ui.ErrPromptCancelled}
		ctx, finalizer := newCtx(t, mainPath, prompts)

		err := runWorktreeRemove(ctx, []string{"feat"}, listFn)

		require.NoError(t, err)
		content, readErr := os.ReadFile(finalizer)
		require.NoError(t, readErr)
		require.Empty(t, content)
	})
}
