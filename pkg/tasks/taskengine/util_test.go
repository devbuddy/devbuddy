package taskengine

import (
	"bytes"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func dummyTask(name string) *taskapi.Task {
	return &taskapi.Task{TaskDefinition: &taskapi.TaskDefinition{Name: name}}
}

type testingAction struct {
	desc            string
	neededResults   []*taskapi.ActionResult
	neededCallCount int

	runResult    error
	runCallCount int
}

func (a *testingAction) Description() string {
	return a.desc
}

func (a *testingAction) Needed(ctx *taskapi.Context) *taskapi.ActionResult {
	result := a.neededResults[a.neededCallCount]
	if result == nil {
		panic("the task should not have been called")
	}
	a.neededCallCount++
	return result
}

func (a *testingAction) Run(ctx *taskapi.Context) error {
	a.runCallCount++
	return a.runResult
}

func newTestingAction(desc string, resultBefore, resultAfter *taskapi.ActionResult, runResult error) *testingAction {
	return &testingAction{
		desc:          desc,
		neededResults: []*taskapi.ActionResult{resultBefore, resultAfter},
		runResult:     runResult,
	}
}

func setupTaskTesting() (*taskapi.Context, *bytes.Buffer) {
	buf, ui := termui.NewTesting(false)

	ctx := &taskapi.Context{
		Project:  project.NewFromPath("/src/myproject"),
		UI:       ui,
		Cfg:      config.NewTestConfig(),
		Env:      env.New([]string{}),
		Features: features.NewFeatureSet(),
	}

	return ctx, buf
}
