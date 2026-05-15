package worktree

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseListPorcelain(t *testing.T) {
	input := `worktree /src/github.com/acme/api
HEAD 1111111111111111111111111111111111111111
branch refs/heads/main

worktree /src/github.com/acme/api--feature-a
HEAD 2222222222222222222222222222222222222222
branch refs/heads/feature-a

worktree /src/github.com/acme/api--detached
HEAD 3333333333333333333333333333333333333333
detached
`

	got, err := ParseListPorcelain(input)

	require.NoError(t, err)
	require.Len(t, got, 3)
	require.Equal(t, "/src/github.com/acme/api", got[0].Path)
	require.Equal(t, "1111111111111111111111111111111111111111", got[0].Head)
	require.Equal(t, "main", got[0].Branch)
	require.False(t, got[0].Detached)
	require.Equal(t, "/src/github.com/acme/api--detached", got[2].Path)
	require.Equal(t, "3333333333333333333333333333333333333333", got[2].Head)
	require.Empty(t, got[2].Branch)
	require.True(t, got[2].Detached)
}

func TestManagedPath(t *testing.T) {
	got, err := ManagedPath("/src/github.com/acme/api", "feature/login flow")

	require.NoError(t, err)
	require.Equal(t, filepath.Clean("/src/github.com/acme/api--feature-login-flow"), got)
}

func TestCheckedOutBranch(t *testing.T) {
	worktrees := []Worktree{
		{Path: "/src/github.com/acme/api", Branch: "main"},
		{Path: "/src/github.com/acme/api--feature-a", Branch: "feature-a"},
	}

	got := CheckedOutBranch(worktrees, "feature-a")

	require.NotNil(t, got)
	require.Equal(t, "/src/github.com/acme/api--feature-a", got.Path)
	require.Nil(t, CheckedOutBranch(worktrees, "missing"))
}
