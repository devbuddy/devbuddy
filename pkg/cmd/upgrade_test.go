package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/updatecheck"
)

func TestRunUpgradeExecutesSharedUpgradePlan(t *testing.T) {
	runner := &upgradeRunnerSpy{}
	plan := updatecheck.Plan{Command: "curl -sSL install.sh | VERSION=v0.17.0 sh"}

	err := runUpgradePlan(runner, plan)

	require.NoError(t, err)
	require.Equal(t, "curl -sSL install.sh | VERSION=v0.17.0 sh", runner.command)
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
