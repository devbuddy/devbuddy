package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/ui"
	"github.com/devbuddy/devbuddy/pkg/worktree"
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
