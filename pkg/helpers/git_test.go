package helpers

import (
	"testing"

	"github.com/Flaque/filet"
	"github.com/pior/dad/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestGitGithubProjectURL(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	test.GitInit(t, tmpdir)

	url, err := NewGitRepo(tmpdir).BuildGithubProjectURL()

	require.NoError(t, err, "BuildGithubProjectURL() failed")
	require.Equal(t, "https://github.com/org1/repo1/tree/master", url)
}

func TestGitGithubPullrequestURL(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	test.GitInit(t, tmpdir)

	url, err := NewGitRepo(tmpdir).BuildGithubPullrequestURL()

	require.NoError(t, err, "BuildGithubProjectURL() failed")
	require.Equal(t, "https://github.com/org1/repo1/pull/master?expand=1", url)
}
