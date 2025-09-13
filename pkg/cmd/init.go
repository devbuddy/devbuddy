package cmd

import (
	"os"
	"slices"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

var initCmd = &cobra.Command{
	Use:          "init [template]",
	Short:        "Initialize a project in the current directory",
	RunE:         initRun,
	Args:         zeroOrOneArg,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func initRun(_ *cobra.Command, args []string) error {
	var templateName string
	if len(args) == 1 {
		templateName = args[0]
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui := termui.New(cfg)

	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	return createManifest(ui, projectPath, templateName)
}

func createManifest(ui *termui.UI, projectPath string, templateName string) error {
	templates := manifest.ListTemplates()

	if templateName == "" || !slices.Contains(templates, templateName) {
		prompt := promptui.Select{
			Label:        "Select a template",
			Items:        templates,
			HideSelected: true,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		templateName = result
	}

	if err := manifest.Create(projectPath, templateName); err != nil {
		return err
	}

	ui.ActionHeader("Created dev.yml with template " + templateName)
	ui.ActionNotice("Open dev.yml to adjust for your needs.")
	return nil
}
