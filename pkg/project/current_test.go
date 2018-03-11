package project

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

func createManifest(dir string) {
	manifestPath := filepath.Join(dir, "dev.yml")
	err := ioutil.WriteFile(manifestPath, []byte("{}"), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func TestFindByPath(t *testing.T) {
	defer filet.CleanUp(t)
	dir := filet.TmpDir(t, "")
	createManifest(dir)

	man, err := findByPath(dir)
	require.NoError(t, err, "findByPath() failed")
	require.NotEqual(t, nil, man)
	require.Equal(t, dir, man.Path)
}

func TestFindByPathNested(t *testing.T) {
	defer filet.CleanUp(t)
	dir := filet.TmpDir(t, "")
	createManifest(dir)

	nestedDir := filepath.Join(dir, "subdir1", "subdir2")
	os.MkdirAll(nestedDir, os.ModePerm)

	man, err := findByPath(nestedDir)
	require.NoError(t, err, "findByPath() failed")
	require.NotEqual(t, nil, man)
	require.Equal(t, dir, man.Path)
}
