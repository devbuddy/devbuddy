package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

var upgradeCmd = &cobra.Command{
	Use:          "upgrade",
	Short:        "[experimental] Upgrade DevBuddy to the latest available release.",
	RunE:         upgradeRun,
	Args:         noArgs,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func upgradeRun(_ *cobra.Command, _ []string) error {
	ctx, err := context.Load(false)
	if err != nil {
		return err
	}

	plateform := fmt.Sprintf("bud-%s-%s", runtime.GOOS, runtime.GOARCH)

	ctx.UI.CommandRun("Getting latest release for", plateform)

	upgrader := helpers.NewUpgrader(ctx, true)
	release, err := upgrader.LatestRelease(plateform)
	if err != nil {
		return err
	}

	destinationPath, err := os.Executable()
	if err != nil {
		return err
	}

	ctx.UI.CommandRun("Downloading", release.DownloadURL)

	if err := upgrader.Perform(destinationPath, release.DownloadURL); err != nil {
		return err
	}

	ctx.UI.CommandActed()
	return nil
}
