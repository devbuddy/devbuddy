package taskengine

import (
	"errors"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/stretchr/testify/require"
)

func TestRunActionNeeded(t *testing.T) {
	ctx, output := setupTaskTesting()
	action := newTestingAction("Action X", api.Needed("some-reason"), api.NotNeeded(), nil)
	task := newTaskWithAction("testtask", action)

	taskRunner := &TaskRunnerImpl{ctx: ctx}
	err := taskRunner.Run(task)
	require.NoError(t, err)

	require.Equal(t, 2, action.neededCallCount)
	require.Equal(t, 1, action.runCallCount)

	require.Contains(t, output.String(), "Action X")
}

func TestRunActionNotNeeded(t *testing.T) {
	ctx, output := setupTaskTesting()
	action := newTestingAction("Action X", api.NotNeeded(), nil, nil)
	task := newTaskWithAction("testtask", action)

	taskRunner := &TaskRunnerImpl{ctx: ctx}
	err := taskRunner.Run(task)
	require.NoError(t, err)

	require.Equal(t, 1, action.neededCallCount)
	require.Equal(t, 0, action.runCallCount)

	require.NotContains(t, output.String(), "Action X")
}

func TestRunActionFailureOnNeeded(t *testing.T) {
	ctx, _ := setupTaskTesting()
	action := newTestingAction("Action X", api.Failed("ERROR_X"), nil, nil)
	task := newTaskWithAction("testtask", action)

	taskRunner := &TaskRunnerImpl{ctx: ctx}
	err := taskRunner.Run(task)
	require.Error(t, err, "failed to detect whether it need to run: ERROR_X")

	require.Equal(t, 1, action.neededCallCount)
	require.Equal(t, 0, action.runCallCount)
}

func TestRunActionFailureOnRun(t *testing.T) {
	ctx, output := setupTaskTesting()
	action := newTestingAction("Action X", api.Needed("some-reason"), nil, errors.New("RunFailed"))
	task := newTaskWithAction("testtask", action)

	taskRunner := &TaskRunnerImpl{ctx: ctx}
	err := taskRunner.Run(task)
	require.EqualError(t, err, `action "Action X": failed to run: RunFailed`)

	require.Equal(t, 1, action.neededCallCount)
	require.Equal(t, 1, action.runCallCount)

	require.Contains(t, output.String(), "Action X")
}

func TestRunActionStillNeeded(t *testing.T) {
	ctx, _ := setupTaskTesting()
	action := newTestingAction("Action X", api.Needed("some-reason"), api.Needed("some-reason"), nil)
	task := newTaskWithAction("testtask", action)

	taskRunner := &TaskRunnerImpl{ctx: ctx}
	err := taskRunner.Run(task)
	require.EqualError(t, err, `action "Action X": ran successfully but still need to run: some-reason`)

	require.Equal(t, 2, action.neededCallCount)
	require.Equal(t, 1, action.runCallCount)
}

func TestTaskRunner(t *testing.T) {
	ctx, output := setupTaskTesting()
	action1 := newTestingAction("Action 1", api.Needed("some-reason"), api.NotNeeded(), nil)
	action2 := newTestingAction("Action 2", api.Needed("some-reason"), api.NotNeeded(), nil)

	task := &api.Task{
		TaskDefinition: &api.TaskDefinition{Name: "testtask"},
		Actions:        []api.TaskAction{action1, action2},
	}

	taskRunner := &TaskRunnerImpl{ctx: ctx}
	err := taskRunner.Run(task)
	require.NoError(t, err)

	require.Equal(t, 2, action1.neededCallCount)
	require.Equal(t, 1, action1.runCallCount)

	require.Equal(t, 2, action2.neededCallCount)
	require.Equal(t, 1, action2.runCallCount)

	require.Contains(t, output.String(), "Action 1")
}

func TestTaskRunnerWithError(t *testing.T) {
	ctx, _ := setupTaskTesting()
	action1 := newTestingAction("Action 1", api.Failed("CRASH 1"), nil, nil)
	action2 := newTestingAction("Action 2", nil, nil, nil)

	task := &api.Task{
		TaskDefinition: &api.TaskDefinition{Name: "testtask"},
		Actions:        []api.TaskAction{action1, action2},
	}

	taskRunner := &TaskRunnerImpl{ctx: ctx}
	err := taskRunner.Run(task)
	require.Error(t, err)
	require.EqualError(t, err, `action "Action 1": detecting whether it needs to run: CRASH 1`)

	require.Equal(t, 1, action1.neededCallCount)
	require.Equal(t, 0, action1.runCallCount)

	require.Equal(t, 0, action2.neededCallCount)
	require.Equal(t, 0, action2.runCallCount)
}
