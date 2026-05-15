package cmd

import (
	"strings"
	"testing"

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
