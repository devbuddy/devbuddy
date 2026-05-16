package tasks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	yaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

func TestRubyOk(t *testing.T) {
	task := ensureLoadTestTask(t, `ruby: 3.3.0`)

	require.Equal(t, "Task Ruby (3.3.0) feature=ruby:3.3.0 actions=3", task.Describe())
	require.Equal(t, "3.3.0", task.Info)
	require.Equal(t, 3, len(task.Actions))
	requireFeature(t, task, "ruby", "3.3.0")
}

func TestRubyMissingVersionNoFile(t *testing.T) {
	_, err := loadTestTask(t, `ruby:`)

	require.Error(t, err, "buildFromDefinition() should have failed without a version")
}

// loadRubyTaskInDir parses a task payload with TaskConfig.ProjectPath set, so
// the parser can consult a .ruby-version file in that directory.
func loadRubyTaskInDir(t *testing.T, payload, projectPath string) (*api.Task, error) {
	t.Helper()
	var data any
	require.NoError(t, yaml.Unmarshal([]byte(payload), &data))

	taskConfig, err := api.NewTaskConfig(data)
	require.NoError(t, err)
	taskConfig.ProjectPath = projectPath

	task := &api.Task{TaskDefinition: api.GetDefinitionOrUnknown("ruby")}
	return task, task.TaskDefinition.Parser(taskConfig, task)
}

func TestRubyVersionFromFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ruby-version"), []byte("3.3.0\n"), 0o600))

	task, err := loadRubyTaskInDir(t, `ruby:`, dir)
	require.NoError(t, err)
	require.Equal(t, "3.3.0", task.Info)
	requireFeature(t, task, "ruby", "3.3.0")
}

func TestRubyVersionFromFileStripsEnginePrefix(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ruby-version"), []byte("ruby-3.3.4\n"), 0o600))

	task, err := loadRubyTaskInDir(t, `ruby:`, dir)
	require.NoError(t, err)
	require.Equal(t, "3.3.4", task.Info)
}

func TestRubyExplicitVersionWinsOverFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ruby-version"), []byte("3.0.0\n"), 0o600))

	task, err := loadRubyTaskInDir(t, `ruby: 3.3.0`, dir)
	require.NoError(t, err)
	require.Equal(t, "3.3.0", task.Info)
}

func TestRubyInvalid(t *testing.T) {
	_, err := loadTestTask(t, `ruby: 3`)

	require.Error(t, err, "buildFromDefinition() should have failed")
}
