package cmd

import (
	"os"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/termui"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project in the current directory",
	Run:   initRun,
	Args:  noArgs,
}

func initRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.NewUI(cfg)

	ui.ActionHeader("Creating a default dev.yml file.")

	projectPath, err := os.Getwd()
	checkError(err)

	err = manifest.Create(projectPath)
	checkError(err)

	ui.ActionNotice("Open dev.yml to adjust for your needs.")
	ui.ActionDone()
}
