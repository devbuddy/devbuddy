package manifest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

var manifestContent = []byte(`
up:
  - task1
  - task2
`)

func createManifest(dir string) {
	manifestPath := filepath.Join(dir, "dev.yml")
	err := ioutil.WriteFile(manifestPath, manifestContent, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func TestLoad(t *testing.T) {
	defer filet.CleanUp(t)

	dir := filet.TmpDir(t, "")
	createManifest(dir)

	man, err := Load(dir)
	require.NoError(t, err, "Load() failed")
	require.NotEqual(t, nil, man)
	require.Equal(t, []interface{}{"task1", "task2"}, man.Up)
}
