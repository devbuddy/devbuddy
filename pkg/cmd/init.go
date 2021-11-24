package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/termui"
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

	ui := termui.New(cfg)

	projectPath, err := os.Getwd()
	checkError(err)

	err = createManifest(ui, projectPath)
	checkError(err)
}

func createManifest(ui *termui.UI, projectPath string) error {
	ui.ActionHeader("Creating a default dev.yml file.")

	err := manifest.Create(projectPath)
	if err != nil {
		return err
	}

	ui.ActionNotice("Open dev.yml to adjust for your needs.")
	return nil
}
