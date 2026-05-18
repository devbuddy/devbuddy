package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/ui"
	"github.com/devbuddy/devbuddy/pkg/updatecheck"
)

func TestRunUpgradePrintsVersionsAsksConfirmationAndExecutesSharedPlan(t *testing.T) {
	runner := &upgradeRunnerSpy{}
	prompts := &ui.FakePrompts{ConfirmValue: true}
	var out bytes.Buffer

	err := runUpgrade(runner, prompts, &out, "v0.16.1 [2026-05-17 12:00:00 +0000 UTC]", updatecheck.Release{Version: "v0.17.0"})

	require.NoError(t, err)
	require.Contains(t, out.String(), "Current version: v0.16.1")
	require.Contains(t, out.String(), "Latest version:  v0.17.0")
	require.Equal(t, []ui.ConfirmRequest{{Label: "Upgrade DevBuddy to v0.17.0 now"}}, prompts.ConfirmRequests)
	require.Equal(t, "curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | VERSION=v0.17.0 sh", runner.command)
}

func TestRunUpgradeSkipsCommandWhenVersionIsCurrent(t *testing.T) {
	runner := &upgradeRunnerSpy{}
	prompts := &ui.FakePrompts{ConfirmValue: true}
	var out bytes.Buffer

	err := runUpgrade(runner, prompts, &out, "v0.17.0", updatecheck.Release{Version: "v0.17.0"})

	require.NoError(t, err)
	require.Contains(t, out.String(), "Current version: v0.17.0")
	require.Contains(t, out.String(), "Latest version:  v0.17.0")
	require.Contains(t, out.String(), "DevBuddy is already at the latest version.")
	require.Empty(t, prompts.ConfirmRequests)
	require.Empty(t, runner.command)
}

func TestRunUpgradeSkipsCommandWhenUserDeclines(t *testing.T) {
	runner := &upgradeRunnerSpy{}
	prompts := &ui.FakePrompts{ConfirmValue: false}
	var out bytes.Buffer

	err := runUpgrade(runner, prompts, &out, "v0.16.1", updatecheck.Release{Version: "v0.17.0"})

	require.NoError(t, err)
	require.Contains(t, out.String(), "Upgrade cancelled.")
	require.Empty(t, runner.command)
}

func TestRunUpgradeReturnsPlanNoteWithoutCommand(t *testing.T) {
	runner := &upgradeRunnerSpy{}
	plan := updatecheck.Plan{Note: "No upgrade command is available."}

	err := runUpgradePlan(runner, plan)

	require.EqualError(t, err, "No upgrade command is available.")
	require.Empty(t, runner.command)
}

func TestCommandVersionUsesRootVersion(t *testing.T) {
	root := &cobra.Command{Use: "bud", Version: "v0.16.1-homebrew"}
	sub := &cobra.Command{Use: "upgrade"}
	root.AddCommand(sub)

	require.Equal(t, "v0.16.1-homebrew", commandVersion(sub))
}

type upgradeRunnerSpy struct {
	command string
	result  *executor.Result
}

func (s *upgradeRunnerSpy) Run(cmd *executor.Command) *executor.Result {
	s.command = cmd.Program
	if s.result != nil {
		return s.result
	}
	return &executor.Result{}
}
