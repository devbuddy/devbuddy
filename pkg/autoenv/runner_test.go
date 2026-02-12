package autoenv

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/autoenv/features"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
	"github.com/devbuddy/devbuddy/pkg/utils"

	"github.com/stretchr/testify/require"
)

type RecorderFeature struct {
	name  string
	calls []string
}

func (r *RecorderFeature) Name() string {
	return r.name
}

func (r *RecorderFeature) Activate(ctx *context.Context, param string) (bool, error) {
	r.calls = append(r.calls, "activate "+param)
	return false, nil
}

func (r *RecorderFeature) Deactivate(ctx *context.Context, param string) {
	r.calls = append(r.calls, "deactivate "+param)
}

func (r *RecorderFeature) getCallsAndReset() []string {
	defer func() { r.calls = nil }()
	return r.calls
}

// WatcherRecorderFeature is a RecorderFeature that also implements FileWatcher.
type WatcherRecorderFeature struct {
	RecorderFeature
	watchedFiles []string
}

func (w *WatcherRecorderFeature) WatchedFiles(_ string) []string {
	return w.watchedFiles
}

func newRunner(env *env.Env, reg *features.Register) *runner {
	return newRunnerWithProject(env, reg, "/project")
}

func newRunnerWithProject(env *env.Env, reg *features.Register, projectPath string) *runner {
	_, ui := termui.NewTesting(false)
	state := &StateManager{env, ui}
	checksums, _ := state.GetFileChecksums()
	return &runner{
		ctx: &context.Context{
			Cfg:     nil,
			Project: project.NewFromPath(projectPath),
			UI:      ui,
			Env:     env,
		},
		state:       state,
		features:    reg,
		fileTracker: utils.NewFileTracker(checksums),
	}
}

func TestRunner(t *testing.T) {
	testEnv := env.New([]string{})

	reg := &features.Register{}
	rust := &RecorderFeature{name: "rust"}
	reg.Register(rust)
	elixir := &RecorderFeature{name: "elixir"}
	reg.Register(elixir)

	runner := newRunner(testEnv, reg)
	runner.sync(NewFeatureSet())
	require.Empty(t, rust.getCallsAndReset())

	// Add a feature
	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(NewFeatureInfo("rust", "1.0")))
	require.Equal(t, []string{"activate 1.0"}, rust.getCallsAndReset())

	// Second time with the same feature
	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(NewFeatureInfo("rust", "1.0")))
	require.Empty(t, rust.getCallsAndReset())

	// Change the feature param
	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(NewFeatureInfo("rust", "2.0")))
	require.Equal(t, []string{"deactivate 1.0", "activate 2.0"}, rust.getCallsAndReset())

	// With no feature
	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet())
	require.Equal(t, []string{"deactivate 2.0"}, rust.getCallsAndReset())
}

func TestRunnerWithTwoFeatures(t *testing.T) {
	testEnv := env.New([]string{})

	reg := &features.Register{}
	rust := &RecorderFeature{name: "rust"}
	reg.Register(rust)
	elixir := &RecorderFeature{name: "elixir"}
	reg.Register(elixir)

	runner := newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(NewFeatureInfo("rust", "1.0")).With(NewFeatureInfo("elixir", "0.4")))
	require.Equal(t, []string{"activate 1.0"}, rust.getCallsAndReset())
	require.Equal(t, []string{"activate 0.4"}, elixir.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(NewFeatureInfo("elixir", "0.4")))
	require.Equal(t, []string{"deactivate 1.0"}, rust.getCallsAndReset())
	require.Empty(t, elixir.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(NewFeatureInfo("rust", "1.0")).With(NewFeatureInfo("elixir", "0.4")))
	require.Equal(t, []string{"activate 1.0"}, rust.getCallsAndReset())
	require.Empty(t, elixir.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet().With(NewFeatureInfo("rust", "1.0")).With(NewFeatureInfo("elixir", "0.5")))
	require.Empty(t, rust.getCallsAndReset())
	require.Equal(t, []string{"deactivate 0.4", "activate 0.5"}, elixir.getCallsAndReset())

	runner = newRunner(testEnv, reg)
	runner.sync(NewFeatureSet())
	require.Equal(t, []string{"deactivate 1.0"}, rust.getCallsAndReset())
	require.Equal(t, []string{"deactivate 0.5"}, elixir.getCallsAndReset())
}

func TestRunnerChangeProject(t *testing.T) {
	testEnv := env.New([]string{})

	reg := &features.Register{}
	rust := &RecorderFeature{name: "rust"}
	reg.Register(rust)

	// Add a feature
	runner := newRunnerWithProject(testEnv, reg, "/project-a")
	runner.sync(NewFeatureSet().With(NewFeatureInfo("rust", "1.0")))
	require.Equal(t, []string{"activate 1.0"}, rust.getCallsAndReset())

	// Same feature on a different project
	runner = newRunnerWithProject(testEnv, reg, "/project-b")
	runner.sync(NewFeatureSet().With(NewFeatureInfo("rust", "1.0")))
	require.Equal(t, []string{"deactivate 1.0", "activate 1.0"}, rust.getCallsAndReset())
}

func TestRunnerFileWatcher(t *testing.T) {
	tmpdir := t.TempDir()
	envPath := filepath.Join(tmpdir, ".env")
	os.WriteFile(envPath, []byte("VAR=one"), 0644)

	testEnv := env.New([]string{})
	reg := &features.Register{}
	watcher := &WatcherRecorderFeature{
		RecorderFeature: RecorderFeature{name: "envfile"},
		watchedFiles:    []string{envPath},
	}
	reg.Register(watcher)

	// Initial activation
	r := newRunner(testEnv, reg)
	r.sync(NewFeatureSet().With(NewFeatureInfo("envfile", envPath)))
	require.Equal(t, []string{"activate " + envPath}, watcher.getCallsAndReset())

	// Same param, file unchanged → no reactivation
	r = newRunner(testEnv, reg)
	r.sync(NewFeatureSet().With(NewFeatureInfo("envfile", envPath)))
	require.Empty(t, watcher.getCallsAndReset())

	// Modify the watched file → reactivation
	os.WriteFile(envPath, []byte("VAR=two"), 0644)

	r = newRunner(testEnv, reg)
	r.sync(NewFeatureSet().With(NewFeatureInfo("envfile", envPath)))
	require.Equal(t, []string{"activate " + envPath}, watcher.getCallsAndReset())

	// File unchanged again → no reactivation
	r = newRunner(testEnv, reg)
	r.sync(NewFeatureSet().With(NewFeatureInfo("envfile", envPath)))
	require.Empty(t, watcher.getCallsAndReset())
}
