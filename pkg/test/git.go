package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// GitInit initializes a path with a simple Git repo
func GitInit(t *testing.T, path string) {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "commit", "-m", "Commit1", "--allow-empty")
	cmd.Dir = path
	cmd.Env = []string{
		"GIT_COMMITTER_NAME=John",
		"GIT_AUTHOR_NAME=John",
		"GIT_COMMITTER_EMAIL=john@doo.com",
		"GIT_AUTHOR_EMAIL=john@doo.com",
	}
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "remote", "add", "origin", "git@github.com:org1/repo1.git")
	cmd.Dir = path
	require.NoError(t, cmd.Run())
}
