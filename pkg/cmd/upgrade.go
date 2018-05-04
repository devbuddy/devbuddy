package cmd

import (
	"os"

	"fmt"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/env"
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
	env := env.NewFromOS()

	cfg, err := config.Load()
	checkError(err)

	u := helpers.NewUpgrader(cfg, true)
	ui := termui.NewUI(cfg)

	ui.CommandHeader(fmt.Sprintf("upgrade %s", env.Platform()))

	release, err := u.LatestRelease()

	destinationPath, err := os.Executable()
	checkError(err)

	ui.CommandRun("Upgrading", destinationPath)

	err = u.Perform(destinationPath, release.DownloadURL)
	checkError(err)

	ui.CommandActed()
}
