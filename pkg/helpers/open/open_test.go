package open

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Flaque/filet"

	"github.com/pior/dad/pkg/manifest"
	"github.com/pior/dad/pkg/project"

	"github.com/stretchr/testify/require"
)

func TestFindLink(t *testing.T) {
	open := map[string]string{"doc": "http://doc.com", "logs": "http://logs"}
	proj := &project.Project{Manifest: &manifest.Manifest{Open: open}}

	_, err := FindLink(proj, "")
	require.Error(t, err)

	_, err = FindLink(proj, "unknown")
	require.Error(t, err)

	url, err := FindLink(proj, "doc")
	require.NoError(t, err)
	require.Equal(t, "http://doc.com", url)
}

func TestFindLinkDefault(t *testing.T) {
	open := map[string]string{"doc": "http://doc.com"}
	proj := &project.Project{Manifest: &manifest.Manifest{Open: open}}

	url, err := FindLink(proj, "")
	require.NoError(t, err)
	require.Equal(t, "http://doc.com", url)
}

func setupProject(t *testing.T, path string) *project.Project {
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "remote", "add", "origin", "git@github.com:org1/repo1.git")
	cmd.Dir = path
	require.NoError(t, cmd.Run())

	return &project.Project{Path: path, Manifest: &manifest.Manifest{}}
}

func TestFindLinkGithub(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	proj := setupProject(t, tmpdir)

	nameToURL := map[string]string{
		"pullrequest": "https://github.com/org1/repo1/pull/master?expand=1",
		"pr":          "https://github.com/org1/repo1/pull/master?expand=1",
		"github":      "https://github.com/org1/repo1/tree/master",
		"gh":          "https://github.com/org1/repo1/tree/master",
	}
	for name, expectedURL := range nameToURL {
		url, err := FindLink(proj, name)
		require.NoError(t, err)
		require.Equal(t, expectedURL, url)
	}
}
