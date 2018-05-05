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

	u := helpers.NewUpgrader(cfg, true)
	ui := termui.NewUI(cfg)

	plateform := fmt.Sprintf("dad-%s-%s", runtime.GOOS, runtime.GOARCH)

	ui.CommandRun("Getting latest release for", plateform)

	release, err := u.LatestRelease(plateform)
	checkError(err)

	destinationPath, err := os.Executable()
	checkError(err)

	ui.CommandRun("Downloading", release.DownloadURL)

	err = u.Perform(destinationPath, release.DownloadURL)
	checkError(err)

	ui.CommandActed()
}
