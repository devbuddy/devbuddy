package cmd

import (
	"testing"

	"github.com/pior/dad/pkg/manifest"
	"github.com/pior/dad/pkg/project"

	"github.com/stretchr/testify/require"
)

func TestFindOpenURL(t *testing.T) {
	open := map[string]string{"doc": "http://doc.com", "logs": "http://logs"}
	proj := &project.Project{Manifest: &manifest.Manifest{Open: open}}

	url, err := findOpenURL(proj, []string{})
	require.Error(t, err)

	url, err = findOpenURL(proj, []string{"unknown"})
	require.Error(t, err)

	url, err = findOpenURL(proj, []string{"doc"})
	require.NoError(t, err)
	require.Equal(t, "http://doc.com", url)
}

func TestFindOpenURLDefault(t *testing.T) {
	open := map[string]string{"doc": "http://doc.com"}
	proj := &project.Project{Manifest: &manifest.Manifest{Open: open}}

	url, err := findOpenURL(proj, []string{})
	require.NoError(t, err)
	require.Equal(t, "http://doc.com", url)
}
