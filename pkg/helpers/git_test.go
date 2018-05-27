package helpers

import (
	"os/exec"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

func buildGitRepos(t *testing.T, path string) {
	init := `
		set -ex
		git init
		git config user.email "you@example.com"
		git config user.name "Your Name"
		git commit -m Commit1 --allow-empty
		git remote add origin git@github.com:org1/repo1.git
	`
	initFile := filet.TmpFile(t, "", init)
	cmd := exec.Command("sh", initFile.Name(), path)
	cmd.Dir = path
	err := cmd.Run()
	require.NoError(t, err, "init failed")
}

func TestGitGithubProjectURL(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	buildGitRepos(t, tmpdir)

	url, err := NewGitRepo(tmpdir).BuildGithubProjectURL()

	require.NoError(t, err, "BuildGithubProjectURL() failed")
	require.Equal(t, "https://github.com/org1/repo1/tree/master", url)
}

func TestGitGithubPullrequestURL(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	buildGitRepos(t, tmpdir)

	url, err := NewGitRepo(tmpdir).BuildGithubPullrequestURL()

	require.NoError(t, err, "BuildGithubProjectURL() failed")
	require.Equal(t, "https://github.com/org1/repo1/pull/master?expand=1", url)
}
