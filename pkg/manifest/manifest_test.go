package manifest

import (
	"testing"

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
  cmd2:
    run: command2
  cmd3: command3

open:
  app: http://localhost:5000
`

func TestLoad(t *testing.T) {
	tmpdir := t.TempDir()

	writer := test.Project(tmpdir)
	writer.Manifest().WriteString(t, manifestContent)

	man, err := Load(tmpdir)
	require.NoError(t, err, "Load() failed")
	require.NotEqual(t, nil, man)

	require.Equal(t, map[string]string{"TESTENV": "TESTVALUE"}, man.Env)
	require.Equal(t, []interface{}{"task1", "task2"}, man.Up)

	commands := map[string]*Command{
		"cmd1": {Run: "command1", Description: "description1"},
		"cmd2": {Run: "command2"},
		"cmd3": {Run: "command3"},
	}
	require.Equal(t, commands, man.GetCommands())

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
