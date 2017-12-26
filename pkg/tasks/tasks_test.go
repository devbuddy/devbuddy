package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func loadTestData(t *testing.T, data string) interface{} {
	var tmp interface{}
	err := yaml.Unmarshal([]byte(data), &tmp)
	if err != nil {
		t.Fatalf("Failed to load a test fixture: %s", err)
	}
	return tmp
}

func TestCustom(t *testing.T) {
	data := loadTestData(t, `
custom:
  met?: test-command
  meet: custom-command
`)

	task, err := BuildFromDefinition(data)
	require.NoError(t, err, "BuildFromDefinition() failed")

	require.Equal(t, task.(*Custom).command, "custom-command")
	require.Equal(t, task.(*Custom).condition, "test-command")
}

func TestPip(t *testing.T) {
	data := loadTestData(t, `
pip:
  - file1
  - file2
`)

	task, err := BuildFromDefinition(data)
	require.NoError(t, err, "BuildFromDefinition() failed")

	require.Equal(t, task.(*Pip).files, []string{"file1", "file2"})
}

func TestPython(t *testing.T) {
	data := loadTestData(t, `python: 3.6.3`)

	task, err := BuildFromDefinition(data)
	require.NoError(t, err, "BuildFromDefinition() failed")

	require.Equal(t, task.(*Python).version, "3.6.3")
}
