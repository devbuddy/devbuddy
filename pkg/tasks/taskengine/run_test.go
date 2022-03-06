package taskengine

import (
	"fmt"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/stretchr/testify/require"
)

type taskRunnerMock struct {
	taskError error
	tasks     []*api.Task
}

func (r *taskRunnerMock) Run(task *api.Task) error {
	r.tasks = append(r.tasks, task)
	return r.taskError
}

type taskSelectorMock struct {
	should bool
}

func (s *taskSelectorMock) ShouldRun(task *api.Task) (bool, error) {
	return s.should, nil
}

func TestRun(t *testing.T) {
	ctx, _ := setupTaskTesting()
	tasks := []*api.Task{dummyTask("1"), dummyTask("2")}

	taskRunner := &taskRunnerMock{}
	taskSelector := &taskSelectorMock{true}

	success, err := Run(ctx, taskRunner, taskSelector, tasks)
	require.NoError(t, err)
	require.True(t, success)

	require.Equal(t, tasks, taskRunner.tasks)
}

func TestRunRequiredTaskCheck(t *testing.T) {
	ctx, _ := setupTaskTesting()
	tasks := []*api.Task{
		&api.Task{TaskDefinition: &api.TaskDefinition{Key: "pip", RequiredTask: "python"}},
	}

	success, err := Run(ctx, nil, nil, tasks)
	require.EqualError(t, err, "You must specify a python task before a pip task")
	require.False(t, success)
}

func TestRunWithTaskError(t *testing.T) {
	ctx, _ := setupTaskTesting()
	tasks := []*api.Task{dummyTask("1"), dummyTask("2")}

	taskRunner := &taskRunnerMock{taskError: fmt.Errorf("oops")}
	taskSelector := &taskSelectorMock{true}

	success, err := Run(ctx, taskRunner, taskSelector, tasks)
	require.NoError(t, err)
	require.False(t, success)

	require.Equal(t, 1, len(taskRunner.tasks))
	require.Equal(t, tasks[0], taskRunner.tasks[0])
}

func TestRunWithTaskWithOsRequirement(t *testing.T) {
	ctx, _ := setupTaskTesting()
	tasks := []*api.Task{dummyTask("1"), dummyTask("2")}

	taskRunner := &taskRunnerMock{taskError: fmt.Errorf("oops")}
	taskSelector := &taskSelectorMock{false}

	success, err := Run(ctx, taskRunner, taskSelector, tasks)
	require.NoError(t, err)
	require.True(t, success)

	require.Equal(t, 0, len(taskRunner.tasks))
}
