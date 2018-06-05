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

commands:
  cmd1:
    desc: description1
    run: command1

open:
  app: http://localhost:5000
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
	require.Equal(t, map[string]*Command{"cmd1": &Command{Run: "command1", Description: "description1"}}, man.Commands)
	require.Equal(t, map[string]string{"app": "http://localhost:5000"}, man.Open)
}
