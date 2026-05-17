package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/updatecheck"
)

var upgradeCmd = &cobra.Command{
	Use:          "upgrade",
	Short:        "[experimental] Upgrade DevBuddy to the latest available release.",
	RunE:         upgradeRun,
	Args:         noArgs,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func upgradeRun(cmd *cobra.Command, _ []string) error {
	ctx, err := context.Load(false)
	if err != nil {
		return err
	}

	ctx.UI.CommandRun("Getting latest DevBuddy release")
	release, err := updatecheck.FetchLatestRelease(nil)
	if err != nil {
		return err
	}

	plan := updatecheck.UpgradePlan(commandVersion(cmd), release.Version)
	return runUpgradePlan(ctx.Executor, plan)
}

func commandVersion(cmd *cobra.Command) string {
	return cmd.Root().Version
}

type upgradeExecutor interface {
	Run(cmd *executor.Command) *executor.Result
}

func runUpgradePlan(runner upgradeExecutor, plan updatecheck.Plan) error {
	if plan.Command == "" {
		if plan.Note != "" {
			return errors.New(plan.Note)
		}
		return errors.New("no upgrade command is available")
	}
	return runner.Run(executor.NewShell(plan.Command)).Error
}
