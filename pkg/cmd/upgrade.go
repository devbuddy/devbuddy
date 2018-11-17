package cmd

import (
	"os"
	"runtime"

	"fmt"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/termui"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "[experimental] Upgrade DevBuddy to the latest available release.",
	Run:   upgradeRun,
	Args:  noArgs,
}

func upgradeRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.New(cfg)

	plateform := fmt.Sprintf("bud-%s-%s", runtime.GOOS, runtime.GOARCH)

	ui.CommandRun("Getting latest release for", plateform)

	upgrader := helpers.NewUpgrader(true)
	release, err := upgrader.LatestRelease(plateform)
	checkError(err)

	destinationPath, err := os.Executable()
	checkError(err)

	ui.CommandRun("Downloading", release.DownloadURL)

	err = upgrader.Perform(ui, destinationPath, release.DownloadURL)
	checkError(err)

	ui.CommandActed()
}
