package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func loadTestTask(t *testing.T, payload string) (*Task, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(payload), &data)
	if err != nil {
		t.Fatalf("Failed to load a test fixture: %s", err)
	}

	task, err := buildFromDefinition(data)
	return task, err
}

func ensureLoadTestTask(t *testing.T, payload string) *Task {
	task, err := loadTestTask(t, payload)
	require.NoError(t, err, "buildFromDefinition() failed")
	return task
}

func dummyTask(name string) *Task {
	return &Task{taskDefinition: &taskDefinition{name: name}}
}

type taskRunnerMock struct {
	taskError error
	tasks     []*Task
}

func (r *taskRunnerMock) Run(ctx *Context, task *Task) error {
	r.tasks = append(r.tasks, task)
	return r.taskError
}

type taskSelectorMock struct {
	should bool
}

func (s *taskSelectorMock) ShouldRun(ctx *Context, task *Task) (bool, error) {
	return s.should, nil
}
