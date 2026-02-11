package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initialize a project in the current directory",
	RunE:         initRun,
	Args:         noArgs,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func initRun(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui := termui.New(cfg)

	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	return createManifest(ui, projectPath)
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
