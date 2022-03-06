package api

import (
	"errors"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/test"

	"github.com/stretchr/testify/require"
)

func newBuilder(description string, runFunc func(*context.Context) error) *taskActionBuilder {
	return &taskActionBuilder{&taskAction{desc: description, runFunc: runFunc}}
}

func TestTaskActionGenericRun(t *testing.T) {
	runCalls := 0

	action := newBuilder("", func(_ *context.Context) error {
		runCalls++
		return nil
	}).taskAction

	action.Run(&context.Context{})
	require.Equal(t, 1, runCalls)
}

func TestTaskActionGenericRunError(t *testing.T) {
	dummyError := errors.New("dummy")

	action := newBuilder("", func(ctx *context.Context) error {
		return dummyError
	}).taskAction

	err := action.Run(&context.Context{})
	require.Equal(t, dummyError, err)
}

func TestTaskActionGenericDescription(t *testing.T) {
	action := newBuilder("dummy desc", nil)

	require.Equal(t, "dummy desc", action.Description())
}

func TestTaskActionGenericFeature(t *testing.T) {
	action := newBuilder("dummy desc", nil)
	action.SetFeature("name", "param")
	require.Equal(t, "name", action.Feature().Name)
	require.Equal(t, "param", action.Feature().Param)
}

func TestTaskActionGenericNoConditions(t *testing.T) {
	action := newBuilder("", func(ctx *context.Context) error { return nil }).taskAction

	result := action.Needed(&context.Context{})
	require.NoError(t, result.Error)
	require.True(t, result.Needed)

	action.Run(&context.Context{})

	result = action.Needed(&context.Context{})
	require.NoError(t, result.Error)
	require.False(t, result.Needed)
}

type recorderCondition struct {
	beforeResult *ActionResult
	beforeCalled int
	afterResult  *ActionResult
	afterCalled  int
}

func (c *recorderCondition) Before(ctx *context.Context) *ActionResult {
	c.beforeCalled++
	return c.beforeResult
}

func (c *recorderCondition) After(ctx *context.Context) *ActionResult {
	c.afterCalled++
	return c.afterResult
}

func TestTaskActionGenericConditions(t *testing.T) {
	cond1 := &recorderCondition{
		beforeResult: NotNeeded(),
		afterResult:  NotNeeded(),
	}
	cond2 := &recorderCondition{
		beforeResult: Needed("pre reason"),
		afterResult:  Needed("post reason"),
	}
	cond3 := &recorderCondition{
		beforeResult: NotNeeded(),
		afterResult:  NotNeeded(),
	}

	builder := newBuilder("", func(ctx *context.Context) error { return nil })
	builder.On(cond1)
	builder.On(cond2)
	builder.On(cond3)
	action := builder.taskAction

	result := action.Needed(&context.Context{})
	require.Equal(t, cond2.beforeResult, result)
	require.Equal(t, 1, cond1.beforeCalled)
	require.Equal(t, 0, cond1.afterCalled)
	require.Equal(t, 1, cond2.beforeCalled)
	require.Equal(t, 0, cond2.afterCalled)
	require.Equal(t, 0, cond3.beforeCalled)
	require.Equal(t, 0, cond3.afterCalled)

	action.Run(&context.Context{})

	result = action.Needed(&context.Context{})
	require.Equal(t, cond2.afterResult, result)
	require.Equal(t, 1, cond1.beforeCalled)
	require.Equal(t, 1, cond1.afterCalled)
	require.Equal(t, 1, cond2.beforeCalled)
	require.Equal(t, 1, cond2.afterCalled)
	require.Equal(t, 0, cond3.beforeCalled)
	require.Equal(t, 0, cond3.afterCalled)
}

func TestTaskActionGenericOnFunc(t *testing.T) {
	calls := 0
	results := []*ActionResult{Needed("reason 1"), NotNeeded()}

	builder := newBuilder("", func(ctx *context.Context) error { return nil })
	builder.On(FuncCondition(func(_ *context.Context) *ActionResult {
		index := calls
		calls++
		return results[index]
	}))
	action := builder.taskAction

	result := action.Needed(&context.Context{})
	require.NoError(t, result.Error)
	require.Equal(t, results[0], result)

	action.Run(&context.Context{})

	result = action.Needed(&context.Context{})
	require.Equal(t, results[1], result)
}

func TestTaskActionGenericFileChange(t *testing.T) {
	tmpdir, tmpfile := test.File(t, "testfile")

	ctx := &context.Context{
		Project: project.NewFromPath(tmpdir),
	}

	runFunc := func(ctx *context.Context) error { return nil }

	// Without file

	action := newBuilder("", runFunc).On(FileCondition("testfile")).taskAction

	result := action.Needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	action.Run(ctx)

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// With a new file

	test.WriteFile(tmpfile, []byte("content-A"))

	action = newBuilder("", runFunc).On(FileCondition("testfile")).taskAction

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.True(t, result.Needed)
	require.Equal(t, "file testfile has changed", result.Reason)

	action.Run(ctx)

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// The file did not change

	action = newBuilder("", runFunc).On(FileCondition("testfile")).taskAction

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// The file changed

	test.WriteFile(tmpfile, []byte("content-B"))

	action = newBuilder("", runFunc).On(FileCondition("testfile")).taskAction

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.True(t, result.Needed)
	require.Equal(t, "file testfile has changed", result.Reason)
}
