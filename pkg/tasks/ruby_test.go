package tasks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/helpers/osidentity"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/ui"
	yaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

func TestRubyOk(t *testing.T) {
	task := ensureLoadTestTask(t, `ruby: 3.3.0`)

	require.Equal(t, "Task Ruby (3.3.0) feature=ruby:3.3.0 actions=3", task.Describe())
	require.Equal(t, "3.3.0", task.Info)
	require.Equal(t, 3, len(task.Actions))
	requireFeature(t, task, "ruby", "3.3.0")
}

func TestRubyMissingVersionNoFile(t *testing.T) {
	_, err := loadTestTask(t, `ruby:`)

	require.Error(t, err, "buildFromDefinition() should have failed without a version")
}

// loadRubyTaskInDir parses a task payload with TaskConfig.ProjectPath set, so
// the parser can consult a .ruby-version file in that directory.
func loadRubyTaskInDir(t *testing.T, payload, projectPath string) (*api.Task, error) {
	t.Helper()
	var data any
	require.NoError(t, yaml.Unmarshal([]byte(payload), &data))

	taskConfig, err := api.NewTaskConfig(data)
	require.NoError(t, err)
	taskConfig.ProjectPath = projectPath

	task := &api.Task{TaskDefinition: api.GetDefinitionOrUnknown("ruby")}
	return task, task.TaskDefinition.Parser(taskConfig, task)
}

func TestRubyVersionFromFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ruby-version"), []byte("3.3.0\n"), 0o600))

	task, err := loadRubyTaskInDir(t, `ruby:`, dir)
	require.NoError(t, err)
	require.Equal(t, "3.3.0", task.Info)
	requireFeature(t, task, "ruby", "3.3.0")
}

func TestRubyVersionFromFileStripsEnginePrefix(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ruby-version"), []byte("ruby-3.3.4\n"), 0o600))

	task, err := loadRubyTaskInDir(t, `ruby:`, dir)
	require.NoError(t, err)
	require.Equal(t, "3.3.4", task.Info)
}

func TestRubyExplicitVersionWinsOverFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ruby-version"), []byte("3.0.0\n"), 0o600))

	task, err := loadRubyTaskInDir(t, `ruby: 3.3.0`, dir)
	require.NoError(t, err)
	require.Equal(t, "3.3.0", task.Info)
}

func TestRubyInvalid(t *testing.T) {
	_, err := loadTestTask(t, `ruby: 3`)

	require.Error(t, err, "buildFromDefinition() should have failed")
}

func TestRubyBundleInstallUsesProjectLocalBundlePath(t *testing.T) {
	task := ensureLoadTestTask(t, `ruby: 3.3.0`)
	runner := &rubyRunner{}
	_, testUI := ui.NewBufferedTesting(false)
	ctx := &context.Context{
		Cfg:     config.NewTestConfig(),
		Env:     env.New([]string{}),
		UI:      testUI,
		Project: project.NewFromPath("/project"),
		Executor: &executor.Executor{
			Runner: runner,
			Cwd:    "/project",
			Env:    env.New([]string{}),
		},
	}

	err := task.Actions[2].Run(ctx)

	require.NoError(t, err)
	require.Equal(t, []string{
		"capture rbenv root",
		"run /rbenv/versions/3.3.0/bin/bundle config set --local path vendor/bundle",
		"run /rbenv/versions/3.3.0/bin/bundle install",
	}, runner.commands)
}

type rubyRunner struct {
	commands []string
}

func (r *rubyRunner) Run(cmd *executor.Command) *executor.Result {
	r.commands = append(r.commands, "run "+rubyCommandString(cmd))
	return &executor.Result{}
}

func (r *rubyRunner) Capture(cmd *executor.Command) *executor.Result {
	r.commands = append(r.commands, "capture "+rubyCommandString(cmd))
	return &executor.Result{Output: "/rbenv\n"}
}

func rubyCommandString(cmd *executor.Command) string {
	result := cmd.Program
	for _, arg := range cmd.Args {
		result += " " + arg
	}
	return result
}

func TestRubyInstallRbenvHomebrewInstallsRubyBuildExplicitly(t *testing.T) {
	ctx, runner := newRubyInstallTestContext(t)

	err := installRbenv(ctx, osidentity.NewMacOSForTest())

	require.NoError(t, err)
	require.Equal(t, []string{"run brew install rbenv ruby-build"}, runner.commands)
	require.Contains(t, runner.runCmds[0].Env, "HOMEBREW_NO_AUTO_UPDATE=1")
}

func TestRubyInstallRbenvDebianUsesAptPackages(t *testing.T) {
	ctx, runner := newRubyInstallTestContext(t)

	err := installRbenv(ctx, osidentity.NewDebianForTest())

	require.NoError(t, err)
	require.Equal(t, []string{
		"run sudo apt-get update",
		"run sudo apt-get install --no-install-recommends -y rbenv ruby-build",
	}, runner.commands)
}

func TestRubyInstallRbenvLinuxFallbackClonesRbenvAndRubyBuild(t *testing.T) {
	ctx, runner := newRubyInstallTestContext(t)

	err := installRbenv(ctx, osidentity.NewLinuxForTest("unknown"))

	require.NoError(t, err)
	require.Equal(t, []string{
		"run git clone https://github.com/rbenv/rbenv.git " + helpers.RbEnvRoot(),
		"run git clone https://github.com/rbenv/ruby-build.git " + filepath.Join(helpers.RbEnvRoot(), "plugins", "ruby-build"),
	}, runner.commands)
}

func TestRubyInstallRbenvLinuxFallbackSkipsExistingCheckouts(t *testing.T) {
	ctx, runner := newRubyInstallTestContext(t)
	require.NoError(t, os.MkdirAll(filepath.Join(helpers.RbEnvRoot(), "plugins", "ruby-build"), 0o755))

	err := installRbenv(ctx, osidentity.NewLinuxForTest("unknown"))

	require.NoError(t, err)
	require.Empty(t, runner.commands)
}

func TestRbenvHasInstallCommand(t *testing.T) {
	ctx, runner := newRubyInstallTestContext(t)
	runner.captureOutputs = []string{helpers.RbEnvRoot(), "global\ninstall\nrehash\n"}

	rbEnv, err := helpers.NewRbEnv(ctx)
	require.NoError(t, err)
	hasRubyBuild, err := rbenvHasInstallCommand(ctx, rbEnv)

	require.NoError(t, err)
	require.True(t, hasRubyBuild)
	require.Equal(t, []string{"capture rbenv root", "capture rbenv commands"}, runner.commands)
}

func TestRbenvHasInstallCommandMissing(t *testing.T) {
	ctx, runner := newRubyInstallTestContext(t)
	runner.captureOutputs = []string{helpers.RbEnvRoot(), "global\nrehash\n"}

	rbEnv, err := helpers.NewRbEnv(ctx)
	require.NoError(t, err)
	hasRubyBuild, err := rbenvHasInstallCommand(ctx, rbEnv)

	require.NoError(t, err)
	require.False(t, hasRubyBuild)
}

func newRubyInstallTestContext(t *testing.T) (*context.Context, *rubyInstallRunner) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("RBENV_ROOT", filepath.Join(os.Getenv("HOME"), ".rbenv"))

	_, testingUI := ui.NewBufferedTesting(false)
	runner := &rubyInstallRunner{}
	ctx := &context.Context{
		Project:  project.NewFromPath("/src/myproject"),
		UI:       testingUI,
		Cfg:      config.NewTestConfig(),
		Env:      env.New([]string{}),
		Executor: &executor.Executor{Runner: runner},
	}
	return ctx, runner
}

type rubyInstallRunner struct {
	commands       []string
	runCmds        []*executor.Command
	captureOutputs []string
	runErr         error
	captureErr     error
}

func (r *rubyInstallRunner) Run(cmd *executor.Command) *executor.Result {
	r.commands = append(r.commands, "run "+formatRubyTestCommand(cmd))
	r.runCmds = append(r.runCmds, cmd)
	if r.runErr != nil {
		return &executor.Result{Error: r.runErr}
	}
	return &executor.Result{}
}

func (r *rubyInstallRunner) Capture(cmd *executor.Command) *executor.Result {
	r.commands = append(r.commands, "capture "+formatRubyTestCommand(cmd))
	if r.captureErr != nil {
		return &executor.Result{Error: r.captureErr}
	}
	output := helpers.RbEnvRoot()
	if len(r.captureOutputs) > 0 {
		output = r.captureOutputs[0]
		r.captureOutputs = r.captureOutputs[1:]
	}
	return &executor.Result{Output: output}
}

func formatRubyTestCommand(cmd *executor.Command) string {
	if cmd.Shell {
		return cmd.Program
	}
	return strings.Join(append([]string{cmd.Program}, cmd.Args...), " ")
}
