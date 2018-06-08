package open

import (
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/test"
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

func TestFindLinkGithub(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	test.GitInit(t, tmpdir)
	proj := &project.Project{Path: tmpdir, Manifest: &manifest.Manifest{}}

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
