package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/ui"
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

	return runUpgrade(ctx.Executor, ctx.UI.Prompts(), cmd.OutOrStdout(), commandVersion(cmd), *release)
}

func commandVersion(cmd *cobra.Command) string {
	return cmd.Root().Version
}

type upgradeExecutor interface {
	Run(cmd *executor.Command) *executor.Result
}

func runUpgrade(runner upgradeExecutor, prompts ui.Prompts, out io.Writer, currentVersion string, release updatecheck.Release) error {
	currentDisplay := displayVersion(currentVersion)
	fmt.Fprintf(out, "Current version: %s\n", currentDisplay)
	fmt.Fprintf(out, "Latest version:  %s\n", release.Version)

	if currentDisplay == release.Version {
		fmt.Fprintln(out, "DevBuddy is already at the latest version.")
		return nil
	}

	confirmed, err := prompts.Confirm(ui.ConfirmRequest{Label: fmt.Sprintf("Upgrade DevBuddy to %s now", release.Version)})
	if errors.Is(err, ui.ErrPromptCancelled) {
		fmt.Fprintln(out, "Upgrade cancelled.")
		return nil
	}
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Fprintln(out, "Upgrade cancelled.")
		return nil
	}

	plan := updatecheck.UpgradePlan(currentVersion, release.Version)
	return runUpgradePlan(runner, plan)
}

func displayVersion(version string) string {
	info := updatecheck.ParseVersion(version)
	if info.Version != "" {
		return info.Version
	}
	return version
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
