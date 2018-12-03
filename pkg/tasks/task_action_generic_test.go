package tasks

import (
	"errors"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/project"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

func TestTaskActionGenericRun(t *testing.T) {
	runCalls := 0

	builder := actionBuilder("", func(ctx *context) error {
		runCalls++
		return nil
	})
	action := builder.Build()

	action.run(&context{})
	require.Equal(t, 1, runCalls)
}

func TestTaskActionGenericRunError(t *testing.T) {
	dummyError := errors.New("dummy")

	builder := actionBuilder("", func(ctx *context) error {
		return dummyError
	})
	action := builder.Build()

	err := action.run(&context{})
	require.Equal(t, dummyError, err)
}

func TestTaskActionGenericDescription(t *testing.T) {
	builder := actionBuilder("dummy desc", nil)
	action := builder.Build()

	require.Equal(t, "dummy desc", action.description())
}

func TestTaskActionGenericNoConditions(t *testing.T) {
	builder := actionBuilder("", func(ctx *context) error { return nil })
	action := builder.Build()

	result := action.needed(&context{})
	require.NoError(t, result.Error)
	require.True(t, result.Needed)

	action.run(&context{})

	result = action.needed(&context{})
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

	result1 := actionNeeded("pre reason")
	result2 := actionNeeded("post reason")

	builder := actionBuilder("", func(ctx *context) error { return nil })
	builder.On(&genericTaskActionCondition{
		pre:  func(ctx *context) *actionResult { pre1Calls++; return actionNotNeeded() },
		post: func(ctx *context) *actionResult { post1Calls++; return actionNotNeeded() },
	})
	builder.On(&genericTaskActionCondition{
		pre:  func(ctx *context) *actionResult { pre2Calls++; return result1 },
		post: func(ctx *context) *actionResult { post2Calls++; return result2 },
	})
	builder.On(&genericTaskActionCondition{
		pre:  func(ctx *context) *actionResult { pre3Calls++; return actionNotNeeded() },
		post: func(ctx *context) *actionResult { post3Calls++; return actionNotNeeded() },
	})
	action := builder.Build()

	result := action.needed(&context{})
	require.Equal(t, result1, result)
	require.Equal(t, 1, pre1Calls)
	require.Equal(t, 0, post1Calls)
	require.Equal(t, 1, pre2Calls)
	require.Equal(t, 0, post2Calls)
	require.Equal(t, 0, pre3Calls)
	require.Equal(t, 0, post3Calls)

	action.run(&context{})

	result = action.needed(&context{})
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
	results := []*actionResult{actionNeeded("reason 1"), actionNotNeeded()}

	builder := actionBuilder("", func(ctx *context) error { return nil })
	builder.OnFunc(func(ctx *context) *actionResult {
		index := calls
		calls++
		return results[index]
	})
	action := builder.Build()

	result := action.needed(&context{})
	require.NoError(t, result.Error)
	require.Equal(t, results[0], result)

	action.run(&context{})

	result = action.needed(&context{})
	require.Equal(t, results[1], result)
}

func TestTaskActionGenericFileChange(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	ctx := &context{
		proj: project.NewFromPath(tmpdir),
	}

	builder := actionBuilder("", func(ctx *context) error { return nil })
	builder.OnFileChange("testfile")

	// Without file

	action := builder.Build()

	result := action.needed(ctx)
	require.NoError(t, result.Error)
	require.True(t, result.Needed)
	require.Equal(t, "file testfile does not exist", result.Reason)

	action.run(ctx)

	result = action.needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// With a new file

	filet.File(t, tmpdir+"/testfile", "content-A")

	action = builder.Build()

	result = action.needed(ctx)
	require.NoError(t, result.Error)
	require.True(t, result.Needed)
	require.Equal(t, "file testfile has changed", result.Reason)

	action.run(ctx)

	result = action.needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// The file did not change

	action = builder.Build()

	result = action.needed(ctx)
	require.NoError(t, result.Error)
	require.False(t, result.Needed)

	// The file changed

	filet.File(t, tmpdir+"/testfile", "content-B")

	action = builder.Build()

	result = action.needed(ctx)
	require.NoError(t, result.Error)
	require.True(t, result.Needed)
	require.Equal(t, "file testfile has changed", result.Reason)

}
