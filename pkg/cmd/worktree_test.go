package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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

func TestWorktreeSwitchMenuTemplateAlignsInactiveRowsWithPandaMarker(t *testing.T) {
	templates := worktreeSwitchTemplates()

	require.Equal(t, "🐼 {{ .Label | cyan }}", templates.Active)
	require.Equal(t, "   {{ .Label }}", templates.Inactive)
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
