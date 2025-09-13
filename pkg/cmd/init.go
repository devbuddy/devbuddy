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
	Use:     "init [template]",
	Short:   "Initialize a project in the current directory",
	Run:     initRun,
	Args:    zeroOrOneArg,
	GroupID: "devbuddy",
}

func initRun(cmd *cobra.Command, args []string) {
	var templateName string
	if len(args) == 1 {
		templateName = args[0]
	}

	cfg, err := config.Load()
	checkError(err)

	ui := termui.New(cfg)

	projectPath, err := os.Getwd()
	checkError(err)

	err = createManifest(ui, projectPath, templateName)
	checkError(err)
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

	ui.ActionHeader("Created dev.yml file with template " + templateName)

	err := manifest.Create(projectPath, templateName)
	if err != nil {
		return err
	}

	ui.ActionNotice("Open dev.yml to adjust for your needs.")
	return nil
}
