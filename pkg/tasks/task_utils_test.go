package tasks

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func loadTestTask(t *testing.T, payload string) (*taskapi.Task, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(payload), &data)
	if err != nil {
		t.Fatalf("Failed to load a test fixture: %s", err)
	}

	task, err := buildFromDefinition(data)
	return task, err
}

func ensureLoadTestTask(t *testing.T, payload string) *taskapi.Task {
	task, err := loadTestTask(t, payload)
	require.NoError(t, err, "buildFromDefinition() failed")
	return task
}

func dummyTask(name string) *taskapi.Task {
	return &Task{TaskDefinition: &TaskDefinition{name: name}}
}

type taskRunnerMock struct {
	taskError error
	tasks     []*taskapi.Task
}

func (r *taskRunnerMock) Run(ctx *taskapi.Context, task *taskapi.Task) error {
	r.tasks = append(r.tasks, task)
	return r.taskError
}

type taskSelectorMock struct {
	should bool
}

func (s *taskSelectorMock) ShouldRun(ctx *taskapi.Context, task *taskapi.Task) (bool, error) {
	return s.should, nil
}
