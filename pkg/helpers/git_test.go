package helpers

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

func buildGitRepos(t *testing.T, path string) {
	init := `
		set -ex

		cd $1
		mkdir repo1

		cd repo1
		git init
		echo POIPOI > file1
		git add file1
		git config user.email "you@example.com"
		git config user.name "Your Name"
		git commit -m Commit1

		cd $1
		git clone repo1 repo2

		cd repo2
		git checkout -b branch1
	`
	initFile := filet.TmpFile(t, "", init)
	cmd := exec.Command("sh", initFile.Name(), path)
	err := cmd.Run()
	require.NoError(t, err, "git repos setup failed")
}

func TestGit(t *testing.T) {
	defer filet.CleanUp(t)

	tmpdir := filet.TmpDir(t, "")
	buildGitRepos(t, tmpdir)

	originPath := filepath.Join(tmpdir, "repo1")
	repoPath := filepath.Join(tmpdir, "repo2")

	gitRepo := NewGitRepo(repoPath)

	branch, err := gitRepo.GetCurrentBranch()
	require.NoError(t, err, "GetCurrentBranch() failed")
	require.Equal(t, "branch1", branch)

	url, err := gitRepo.GetRemoteURL()
	require.NoError(t, err, "GetRemoteURL() failed")
	require.Equal(t, originPath, url)
}
