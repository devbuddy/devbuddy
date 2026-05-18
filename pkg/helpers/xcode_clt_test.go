package helpers

import (
	"errors"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/helpers/osidentity"
	"github.com/devbuddy/devbuddy/pkg/ui"
	"github.com/stretchr/testify/require"
)

func TestEnsureXcodeCommandLineToolsNoopsOutsideMacOS(t *testing.T) {
	runner := &xcodeCLTRunner{}
	ctx := newXcodeCLTContext(runner)

	err := ensureXcodeCommandLineTools(ctx, osidentity.NewDebianForTest(), false)

	require.NoError(t, err)
	require.Empty(t, runner.commands)
}

func TestEnsureXcodeCommandLineToolsReturnsNilWhenInstalled(t *testing.T) {
	runner := &xcodeCLTRunner{}
	ctx := newXcodeCLTContext(runner)

	err := ensureXcodeCommandLineTools(ctx, osidentity.NewMacOSForTest(), false)

	require.NoError(t, err)
	require.Equal(t, []string{"capture xcode-select -p", "capture xcodebuild -checkFirstLaunchStatus"}, runner.commands)
}

func TestEnsureXcodeCommandLineToolsStartsInstallerWhenMissing(t *testing.T) {
	runner := &xcodeCLTRunner{
		xcodeSelectErr: errors.New("missing"),
	}
	ctx := newXcodeCLTContext(runner)

	err := ensureXcodeCommandLineTools(ctx, osidentity.NewMacOSForTest(), false)

	require.EqualError(t, err, "Xcode Command Line Tools are required. Complete the installer dialog, then re-run bud up.")
	require.Equal(t, []string{"capture xcode-select -p", "run xcode-select --install"}, runner.commands)
}

func TestEnsureXcodeCommandLineToolsReportsInstallerLaunchFailure(t *testing.T) {
	runner := &xcodeCLTRunner{
		xcodeSelectErr: errors.New("missing"),
		installerErr:   errors.New("launch failed"),
	}
	ctx := newXcodeCLTContext(runner)

	err := ensureXcodeCommandLineTools(ctx, osidentity.NewMacOSForTest(), false)

	require.EqualError(t, err, "failed to start Xcode Command Line Tools installer: launch failed")
	require.Equal(t, []string{"capture xcode-select -p", "run xcode-select --install"}, runner.commands)
}

func TestEnsureXcodeCommandLineToolsReportsFirstLaunchWhenNonInteractive(t *testing.T) {
	runner := &xcodeCLTRunner{
		xcodeBuildErr: errors.New("first launch incomplete"),
	}
	ctx := newXcodeCLTContext(runner)

	err := ensureXcodeCommandLineTools(ctx, osidentity.NewMacOSForTest(), false)

	require.EqualError(t, err, "Xcode first-launch setup is required. Run `sudo xcodebuild -runFirstLaunch`, then re-run bud up.")
	require.Equal(t, []string{"capture xcode-select -p", "capture xcodebuild -checkFirstLaunchStatus"}, runner.commands)
}

func TestEnsureXcodeCommandLineToolsRunsFirstLaunchWhenInteractive(t *testing.T) {
	runner := &xcodeCLTRunner{
		xcodeBuildErr: errors.New("first launch incomplete"),
	}
	ctx := newXcodeCLTContext(runner)

	err := ensureXcodeCommandLineTools(ctx, osidentity.NewMacOSForTest(), true)

	require.NoError(t, err)
	require.Equal(t, []string{
		"capture xcode-select -p",
		"capture xcodebuild -checkFirstLaunchStatus",
		"run sudo xcodebuild -runFirstLaunch",
	}, runner.commands)
}

func TestEnsureXcodeCommandLineToolsReportsFirstLaunchFailure(t *testing.T) {
	runner := &xcodeCLTRunner{
		xcodeBuildErr:     errors.New("first launch incomplete"),
		firstLaunchRunErr: errors.New("sudo failed"),
	}
	ctx := newXcodeCLTContext(runner)

	err := ensureXcodeCommandLineTools(ctx, osidentity.NewMacOSForTest(), true)

	require.EqualError(t, err, "failed to run Xcode first-launch setup: sudo failed")
	require.Equal(t, []string{
		"capture xcode-select -p",
		"capture xcodebuild -checkFirstLaunchStatus",
		"run sudo xcodebuild -runFirstLaunch",
	}, runner.commands)
}

func newXcodeCLTContext(runner *xcodeCLTRunner) *context.Context {
	_, testUI := ui.NewBufferedTesting(false)
	return &context.Context{
		Cfg: config.NewTestConfig(),
		Env: env.New([]string{}),
		UI:  testUI,
		Executor: &executor.Executor{
			Runner: runner,
			Env:    env.New([]string{}),
		},
	}
}

type xcodeCLTRunner struct {
	commands          []string
	xcodeSelectErr    error
	xcodeBuildErr     error
	installerErr      error
	firstLaunchRunErr error
}

func (r *xcodeCLTRunner) Run(cmd *executor.Command) *executor.Result {
	r.commands = append(r.commands, "run "+commandString(cmd))
	if cmd.Program == "sudo" {
		return &executor.Result{Error: r.firstLaunchRunErr}
	}
	return &executor.Result{Error: r.installerErr}
}

func (r *xcodeCLTRunner) Capture(cmd *executor.Command) *executor.Result {
	r.commands = append(r.commands, "capture "+commandString(cmd))
	if cmd.Program == "xcodebuild" {
		return &executor.Result{Error: r.xcodeBuildErr}
	}
	return &executor.Result{Error: r.xcodeSelectErr}
}

func commandString(cmd *executor.Command) string {
	result := cmd.Program
	for _, arg := range cmd.Args {
		result += " " + arg
	}
	return result
}
