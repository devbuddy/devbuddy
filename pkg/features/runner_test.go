package features

import (
	"strings"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/termui"
	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/config"
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

func newMockEnv(name string, reg *featureRegister, rec *recorder) {
	reg.register(
		name,
		func(param string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
			rec.record("activate", param)
			return false, nil
		},
		func(param string, cfg *config.Config, env *env.Env) {
			rec.record("deactivate", param)
		},
	)
}

func newRunner(env *env.Env, reg *featureRegister) *runner {
	_, ui := termui.NewTesting(false)
	return &runner{cfg: nil, proj: nil, ui: ui, env: env, reg: reg}
}

func TestRunner(t *testing.T) {
	testEnv := env.New([]string{})

	reg := newFeatureRegister()

	rustCalls := newRecorder()
	newMockEnv("rust", reg, rustCalls)

	elixirCalls := newRecorder()
	newMockEnv("elixir", reg, elixirCalls)

	runner := newRunner(testEnv, reg)
	runner.sync(map[string]string{})
	require.Equal(t, []string{}, rustCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(map[string]string{"rust": "1.0"})
	require.Equal(t, []string{"activate 1.0"}, rustCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(map[string]string{"rust": "2.0"})
	require.Equal(t, []string{"deactivate 1.0", "activate 2.0"}, rustCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(map[string]string{})
	require.Equal(t, []string{"deactivate 2.0"}, rustCalls.getCallsAndReset())
}

func TestRunnerWithTwoFeatures(t *testing.T) {
	testEnv := env.New([]string{})

	reg := newFeatureRegister()

	rustCalls := newRecorder()
	newMockEnv("rust", reg, rustCalls)

	elixirCalls := newRecorder()
	newMockEnv("elixir", reg, elixirCalls)

	runner := newRunner(testEnv, reg)
	runner.sync(map[string]string{"rust": "1.0", "elixir": "0.4"})
	require.Equal(t, []string{"activate 1.0"}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{"activate 0.4"}, elixirCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(map[string]string{"elixir": "0.4"})
	require.Equal(t, []string{"deactivate 1.0"}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{}, elixirCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(map[string]string{"rust": "1.0", "elixir": "0.4"})
	require.Equal(t, []string{"activate 1.0"}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{}, elixirCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(map[string]string{"rust": "1.0", "elixir": "0.5"})
	require.Equal(t, []string{}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{"deactivate 0.4", "activate 0.5"}, elixirCalls.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(map[string]string{})
	require.Equal(t, []string{"deactivate 1.0"}, rustCalls.getCallsAndReset())
	require.Equal(t, []string{"deactivate 0.5"}, elixirCalls.getCallsAndReset())
}
