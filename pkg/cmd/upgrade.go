package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/termui"
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
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui := termui.New(cfg)

	plateform := fmt.Sprintf("bud-%s-%s", runtime.GOOS, runtime.GOARCH)

	ui.CommandRun("Getting latest release for", plateform)

	upgrader := helpers.NewUpgrader(true)
	release, err := upgrader.LatestRelease(plateform)
	if err != nil {
		return err
	}

	destinationPath, err := os.Executable()
	if err != nil {
		return err
	}

	ui.CommandRun("Downloading", release.DownloadURL)

	if err := upgrader.Perform(ui, destinationPath, release.DownloadURL); err != nil {
		return err
	}

	ui.CommandActed()
	return nil
}
