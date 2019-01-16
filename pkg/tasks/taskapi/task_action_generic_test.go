package taskapi

import (
	"errors"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/project"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

func newBuilder(description string, runFunc func(*Context) error) *genericTaskActionBuilder {
	return &genericTaskActionBuilder{&genericTaskAction{desc: description, runFunc: runFunc}}
}

func TestTaskActionGenericRun(t *testing.T) {
	runCalls := 0

	action := newBuilder("", func(ctx *Context) error {
		runCalls++
		return nil
	}).genericTaskAction

	action.Run(&Context{})
	require.Equal(t, 1, runCalls)
}

func TestTaskActionGenericRunError(t *testing.T) {
	dummyError := errors.New("dummy")

	action := newBuilder("", func(ctx *Context) error {
		return dummyError
	}).genericTaskAction

	err := action.Run(&Context{})
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
	action := newBuilder("", func(ctx *Context) error { return nil }).genericTaskAction

	result := action.Needed(&Context{})
	require.NoError(t, result.Error)
	require.True(t, result.Needed)

	action.Run(&Context{})

	result = action.Needed(&Context{})
	require.NoError(t, result.Error)
	require.False(t, result.Needed)
}

func TestTaskActionGenericConditions(t *testing.T) {
	pre1Calls := 0
	post1Calls := 0
	pre2Calls := 0
	post2Calls := 0
	pre3Calls := 0
	post3Calls := 0

	result1 := ActionNeeded("pre reason")
	result2 := ActionNeeded("post reason")

	builder := newBuilder("", func(ctx *Context) error { return nil })
	builder.On(&genericTaskActionCondition{
		pre:  func(ctx *Context) *ActionResult { pre1Calls++; return ActionNotNeeded() },
		post: func(ctx *Context) *ActionResult { post1Calls++; return ActionNotNeeded() },
	})
	builder.On(&genericTaskActionCondition{
		pre:  func(ctx *Context) *ActionResult { pre2Calls++; return result1 },
		post: func(ctx *Context) *ActionResult { post2Calls++; return result2 },
	})
	builder.On(&genericTaskActionCondition{
		pre:  func(ctx *Context) *ActionResult { pre3Calls++; return ActionNotNeeded() },
		post: func(ctx *Context) *ActionResult { post3Calls++; return ActionNotNeeded() },
	})
	action := builder.genericTaskAction

	result := action.Needed(&Context{})
	require.Equal(t, result1, result)
	require.Equal(t, 1, pre1Calls)
	require.Equal(t, 0, post1Calls)
	require.Equal(t, 1, pre2Calls)
	require.Equal(t, 0, post2Calls)
	require.Equal(t, 0, pre3Calls)
	require.Equal(t, 0, post3Calls)

	action.Run(&Context{})

	result = action.Needed(&Context{})
	require.Equal(t, result2, result)
	require.Equal(t, 1, pre1Calls)
	require.Equal(t, 1, post1Calls)
	require.Equal(t, 1, pre2Calls)
	require.Equal(t, 1, post2Calls)
	require.Equal(t, 0, pre3Calls)
	require.Equal(t, 0, post3Calls)
}

func TestTaskActionGenericOnFunc(t *testing.T) {
	calls := 0
	results := []*ActionResult{ActionNeeded("reason 1"), ActionNotNeeded()}

	builder := newBuilder("", func(ctx *Context) error { return nil })
	builder.OnFunc(func(ctx *Context) *ActionResult {
		index := calls
		calls++
		return results[index]
	})
	action := builder.genericTaskAction

	result := action.Needed(&Context{})
	require.NoError(t, result.Error)
	require.Equal(t, results[0], result)

	action.Run(&Context{})

	result = action.Needed(&Context{})
	require.Equal(t, results[1], result)
}

func TestTaskActionGenericFileChange(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	ctx := &Context{
		Project: project.NewFromPath(tmpdir),
	}

	runFunc := func(ctx *Context) error { return nil }

	// Without file

	action := newBuilder("", runFunc).OnFileChange("testfile").genericTaskAction

	result := action.Needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	action.Run(ctx)

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// With a new file

	filet.File(t, tmpdir+"/testfile", "content-A")

	action = newBuilder("", runFunc).OnFileChange("testfile").genericTaskAction

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.True(t, result.Needed)
	require.Equal(t, "file testfile has changed", result.Reason)

	action.Run(ctx)

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// The file did not change

	action = newBuilder("", runFunc).OnFileChange("testfile").genericTaskAction

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// The file changed

	filet.File(t, tmpdir+"/testfile", "content-B")

	action = newBuilder("", runFunc).OnFileChange("testfile").genericTaskAction

	result = action.Needed(ctx)
	require.NoError(t, result.Error)
	require.True(t, result.Needed)
	require.Equal(t, "file testfile has changed", result.Reason)

}
