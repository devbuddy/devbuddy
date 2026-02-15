package taskengine

import (
	"bytes"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func dummyTask(name string) *api.Task {
	return &api.Task{TaskDefinition: &api.TaskDefinition{Name: name}}
}

type testingAction struct {
	desc            string
	neededResults   []*api.ActionResult
	neededCallCount int

	feature *autoenv.FeatureInfo

	runResult    error
	runCallCount int
}

func (a *testingAction) Description() string {
	return a.desc
}

func (a *testingAction) Needed(ctx *context.Context) *api.ActionResult {
	result := a.neededResults[a.neededCallCount]
	if result == nil {
		panic("the task should not have been called")
	}
	a.neededCallCount++
	return result
}

func (a *testingAction) Run(ctx *context.Context) error {
	a.runCallCount++
	return a.runResult
}

func (a *testingAction) Feature() *autoenv.FeatureInfo {
	return a.feature
}

func newTestingAction(desc string, resultBefore, resultAfter *api.ActionResult, runResult error) *testingAction {
	return &testingAction{
		desc:          desc,
		neededResults: []*api.ActionResult{resultBefore, resultAfter},
		runResult:     runResult,
	}
}

func newTaskWithAction(name string, action api.TaskAction) *api.Task {
	return &api.Task{
		TaskDefinition: &api.TaskDefinition{Name: name},
		Actions:        []api.TaskAction{action},
	}
}

func setupTaskTesting() (*context.Context, *bytes.Buffer) {
	buf, ui := termui.NewTesting(false)

	ctx := &context.Context{
		Project:  project.NewFromPath("/src/myproject"),
		UI:       ui,
		Cfg:      config.NewTestConfig(),
		Env:      env.New([]string{}),
		Executor: executor.NewExecutor(),
	}

	return ctx, buf
}
