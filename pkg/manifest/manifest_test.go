package manifest

import (
	"testing"

	"github.com/Flaque/filet"
	"github.com/devbuddy/devbuddy/pkg/test"
	"github.com/stretchr/testify/require"
)

var manifestContent = `
env:
  TESTENV: TESTVALUE

up:
  - task1
  - task2

commands:
  cmd1:
    desc: description1
    run: command1

open:
  app: http://localhost:5000
`

func TestLoad(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)
	writer := test.Project(tmpdir)
	writer.Manifest().WriteString(t, manifestContent)

	man, err := Load(tmpdir)
	require.NoError(t, err, "Load() failed")
	require.NotEqual(t, nil, man)

	require.Equal(t, map[string]string{"TESTENV": "TESTVALUE"}, man.Env)
	require.Equal(t, []interface{}{"task1", "task2"}, man.Up)
	require.Equal(t, map[string]*Command{"cmd1": {Run: "command1", Description: "description1"}}, man.Commands)
	require.Equal(t, map[string]string{"app": "http://localhost:5000"}, man.Open)
}

func TestLoadErr(t *testing.T) {
	man, err := Load("")
	require.Error(t, err)
	require.Nil(t, man)

	man, err = Load("/dev/null")
	require.Error(t, err)
	require.Nil(t, man)
}
