package autoenv

import (
	"strings"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/autoenv/features"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/termui"
	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type recorder struct {
	entries []string
}

func newRecorder() *recorder {
	return &recorder{entries: []string{}}
}

func (r *recorder) record(s ...string) {
	r.entries = append(r.entries, strings.Join(s, " "))
}

func (r *recorder) getCallsAndReset() []string {
	val := r.entries
	r.entries = []string{}
	return val
}

func newMockEnv(name string, reg *features.MutableRegister, rec *recorder) {
	reg.Register(
		name,
		func(ctx *context.Context, param string) (bool, error) {
			rec.record("activate", param)
			return false, nil
		},
		func(ctx *context.Context, param string) {
			rec.record("deactivate", param)
		},
	)
}

func newRunner(env *env.Env, reg *features.MutableRegister) *runner {
	return newRunnerWithProject(env, reg, "/project")
}

func newRunnerWithProject(env *env.Env, reg *features.MutableRegister, projectPath string) *runner {
	_, ui := termui.NewTesting(false)
	return &runner{
		ctx: &context.Context{
			Cfg:     nil,
			Project: project.NewFromPath(projectPath),
			UI:      ui,
			Env:     env,
		},
		state: &FeatureState{env, ui},
		reg:   reg,
	}
}

func TestRunner(t *testing.T) {
	testEnv := env.New([]string{})

	reg := features.NewRegister()

	rustCalls := newRecorder()
	newMockEnv("rust", reg, rustCalls)

	elixirCalls := newRecorder()
	newMockEnv("elixir", reg, elixirCalls)

	runner := newRunner(testEnv, reg)
	runner.sync(NewFeatureSet())
	require.Equal(t, []string{}, rustCalls.getCallsAndReset())

	// Add a feature
	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(&FeatureInfo{"rust", "1.0"}))
	require.Equal(t, []string{"activate 1.0"}, rustCalls.getCallsAndReset())

	// Second time with the same feature
	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(&FeatureInfo{"rust", "1.0"}))
	require.Equal(t, []string{}, rustCalls.getCallsAndReset())

	// Change the feature param
	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(&FeatureInfo{"rust", "2.0"}))
	require.Equal(t, []string{"deactivate 1.0", "activate 2.0"}, rustCalls.getCallsAndReset())

	// With no feature
	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet())
	require.Equal(t, []string{"deactivate 2.0"}, rustCalls.getCallsAndReset())
}

func TestRunnerWithTwoFeatures(t *testing.T) {
	testEnv := env.New([]string{})

	reg := features.NewRegister()

	rustCalls := newRecorder()
	newMockEnv("rust", reg, rustCalls)

	elixirCalls := newRecorder()
	newMockEnv("elixir", reg, elixirCalls)

	runner := newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(&FeatureInfo{"rust", "1.0"}).With(&FeatureInfo{"elixir", "0.4"}))
	require.Equal(t, []string{"activate 1.0"}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{"activate 0.4"}, elixirCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(&FeatureInfo{"elixir", "0.4"}))
	require.Equal(t, []string{"deactivate 1.0"}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{}, elixirCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(&FeatureInfo{"rust", "1.0"}).With(&FeatureInfo{"elixir", "0.4"}))
	require.Equal(t, []string{"activate 1.0"}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{}, elixirCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(&FeatureInfo{"rust", "1.0"}).With(&FeatureInfo{"elixir", "0.5"}))
	require.Equal(t, []string{}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{"deactivate 0.4", "activate 0.5"}, elixirCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet())
	require.Equal(t, []string{"deactivate 1.0"}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{"deactivate 0.5"}, elixirCalls.getCallsAndReset())
}

func TestRunnerChangeProject(t *testing.T) {
	testEnv := env.New([]string{})

	reg := features.NewRegister()

	rustCalls := newRecorder()
	newMockEnv("rust", reg, rustCalls)

	// Add a feature
	runner := newRunnerWithProject(testEnv, reg, "/project-a")
	runner.sync(NewFeatureSet().With(&FeatureInfo{"rust", "1.0"}))
	require.Equal(t, []string{"activate 1.0"}, rustCalls.getCallsAndReset())

	// Same feature on a different project
	runner = newRunnerWithProject(testEnv, reg, "/project-b")
	runner.sync(NewFeatureSet().With(&FeatureInfo{"rust", "1.0"}))
	require.Equal(t, []string{"deactivate 1.0", "activate 1.0"}, rustCalls.getCallsAndReset())
}
