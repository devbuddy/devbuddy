package taskapi

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"

	"github.com/stretchr/testify/require"
)

func dummyTask(name string) *Task {
	return &Task{TaskDefinition: &TaskDefinition{Name: name}}
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

type testingAction struct {
	desc            string
	neededResults   []*ActionResult
	neededCallCount int

	runResult    error
	runCallCount int
}

func (a *testingAction) Description() string {
	return a.desc
}

func (a *testingAction) Needed(ctx *Context) *ActionResult {
	result := a.neededResults[a.neededCallCount]
	if result == nil {
		panic("the task should not have been called")
	}
	a.neededCallCount++
	return result
}

func (a *testingAction) Run(ctx *Context) error {
	a.runCallCount++
	return a.runResult
}

func newTestingAction(desc string, resultBefore, resultAfter *ActionResult, runResult error) *testingAction {
	return &testingAction{
		desc:          desc,
		neededResults: []*ActionResult{resultBefore, resultAfter},
		runResult:     runResult,
	}
}

func setupTaskTesting() (*Context, *bytes.Buffer) {
	buf, ui := termui.NewTesting(false)

	ctx := &Context{
		Project:  project.NewFromPath("/src/myproject"),
		UI:       ui,
		Cfg:      config.NewTestConfig(),
		Env:      env.New([]string{}),
		Features: features.NewFeatureSet(),
	}

	return ctx, buf
}

func TestRunActionNeeded(t *testing.T) {
	ctx, output := setupTaskTesting()
	action := newTestingAction("Action X", ActionNeeded("some-reason"), ActionNotNeeded(), nil)

	err := runAction(ctx, action)
	require.NoError(t, err)

	require.Equal(t, 2, action.neededCallCount)
	require.Equal(t, 1, action.runCallCount)

	require.Contains(t, output.String(), "Action X")
}

func TestRunActionNotNeeded(t *testing.T) {
	ctx, output := setupTaskTesting()
	action := newTestingAction("Action X", ActionNotNeeded(), nil, nil)

	err := runAction(ctx, action)
	require.NoError(t, err)

	require.Equal(t, 1, action.neededCallCount)
	require.Equal(t, 0, action.runCallCount)

	require.NotContains(t, output.String(), "Action X")
}

func TestRunActionFailureOnNeeded(t *testing.T) {
	ctx, _ := setupTaskTesting()
	action := newTestingAction("Action X", ActionFailed("ERROR_X"), nil, nil)

	err := runAction(ctx, action)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to detect whether it need to run: ERROR_X")

	require.Equal(t, 1, action.neededCallCount)
	require.Equal(t, 0, action.runCallCount)
}

func TestRunActionFailureOnRun(t *testing.T) {
	ctx, output := setupTaskTesting()
	action := newTestingAction("Action X", ActionNeeded("some-reason"), nil, errors.New("RunFailed"))

	err := runAction(ctx, action)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to run: RunFailed")

	require.Equal(t, 1, action.neededCallCount)
	require.Equal(t, 1, action.runCallCount)

	require.Contains(t, output.String(), "Action X")
}

func TestRunActionStillNeeded(t *testing.T) {
	ctx, _ := setupTaskTesting()
	action := newTestingAction("Action X", ActionNeeded("some-reason"), ActionNeeded("some-reason"), nil)

	err := runAction(ctx, action)
	require.Error(t, err)
	require.Contains(t, err.Error(), "did not produce the expected result: some-reason")

	require.Equal(t, 2, action.neededCallCount)
	require.Equal(t, 1, action.runCallCount)
}

func TestTaskRunner(t *testing.T) {
	ctx, output := setupTaskTesting()
	action1 := newTestingAction("Action 1", ActionNeeded("some-reason"), ActionNotNeeded(), nil)
	action2 := newTestingAction("Action 2", ActionNeeded("some-reason"), ActionNotNeeded(), nil)

	task := &Task{
		TaskDefinition: &TaskDefinition{Name: "testtask"},
		actions:        []TaskAction{action1, action2},
	}

	taskRunner := &TaskRunnerImpl{}
	err := taskRunner.Run(ctx, task)
	require.NoError(t, err)

	require.Equal(t, 2, action1.neededCallCount)
	require.Equal(t, 1, action1.runCallCount)

	require.Equal(t, 2, action2.neededCallCount)
	require.Equal(t, 1, action2.runCallCount)

	require.Contains(t, output.String(), "Action 1")
}

func TestTaskRunnerWithError(t *testing.T) {
	ctx, _ := setupTaskTesting()
	action1 := newTestingAction("Action 1", ActionFailed("CRASH 1"), nil, nil)
	action2 := newTestingAction("Action 2", nil, nil, nil)

	task := &Task{
		TaskDefinition: &TaskDefinition{Name: "testtask"},
		actions:        []TaskAction{action1, action2},
	}

	taskRunner := &TaskRunnerImpl{}
	err := taskRunner.Run(ctx, task)
	require.Error(t, err)
	require.Contains(t, err.Error(), "The task action (Action 1) failed to detect whether it need to run: CRASH 1")

	require.Equal(t, 1, action1.neededCallCount)
	require.Equal(t, 0, action1.runCallCount)

	require.Equal(t, 0, action2.neededCallCount)
	require.Equal(t, 0, action2.runCallCount)
}

func TestRun(t *testing.T) {
	ctx, _ := setupTaskTesting()
	tasks := []*Task{dummyTask("1"), dummyTask("2")}

	taskRunner := &taskRunnerMock{}
	taskSelector := &taskSelectorMock{true}

	success, err := Run(ctx, taskRunner, taskSelector, tasks)
	require.NoError(t, err)
	require.True(t, success)

	require.Equal(t, tasks, taskRunner.tasks)
}

func TestRunRequiredTaskCheck(t *testing.T) {
	ctx, _ := setupTaskTesting()
	tasks := []*Task{
		&Task{TaskDefinition: &TaskDefinition{Key: "pip", RequiredTask: "python"}},
	}

	success, err := Run(ctx, nil, nil, tasks)
	require.EqualError(t, err, "You must specify a python task before a pip task")
	require.False(t, success)
}

func TestRunWithTaskError(t *testing.T) {
	ctx, _ := setupTaskTesting()
	tasks := []*Task{dummyTask("1"), dummyTask("2")}

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
	tasks := []*Task{dummyTask("1"), dummyTask("2")}

	taskRunner := &taskRunnerMock{taskError: fmt.Errorf("oops")}
	taskSelector := &taskSelectorMock{false}

	success, err := Run(ctx, taskRunner, taskSelector, tasks)
	require.NoError(t, err)
	require.True(t, success)

	require.Equal(t, 0, len(taskRunner.tasks))
}
