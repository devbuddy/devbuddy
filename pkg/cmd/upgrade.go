package cmd

import (
	"os"
	"runtime"

	"fmt"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/termui"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "[experimental] Upgrade dad to the latest available release.",
	Run:   upgradeRun,
	Args:  noArgs,
}

func upgradeRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.NewUI(cfg)

	plateform := fmt.Sprintf("dad-%s-%s", runtime.GOOS, runtime.GOARCH)

	ui.CommandRun("Getting latest release for", plateform)

	upgrader := helpers.NewUpgrader(cfg, true)
	release, err := upgrader.LatestRelease(plateform)
	checkError(err)

	destinationPath, err := os.Executable()
	checkError(err)

	ui.CommandRun("Downloading", release.DownloadURL)

	err = upgrader.Perform(destinationPath, release.DownloadURL)
	checkError(err)

	ui.CommandActed()
}
