package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/devbuddy/devbuddy/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestFindByPath(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)
	writer := test.Project(tmpdir)
	writer.Manifest().Empty(t)

	proj, err := findByPath(tmpdir)
	require.NoError(t, err, "findByPath() failed")
	require.NotEqual(t, nil, proj)
	require.Equal(t, tmpdir, proj.Path)
}

func TestFindByPathNested(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)
	writer := test.Project(tmpdir)
	writer.Manifest().Empty(t)

	nestedDir := filepath.Join(tmpdir, "subdir1", "subdir2")
	os.MkdirAll(nestedDir, os.ModePerm)

	proj, err := findByPath(nestedDir)
	require.NoError(t, err, "findByPath() failed")
	require.NotEqual(t, nil, proj)
	require.Equal(t, tmpdir, proj.Path)
}
