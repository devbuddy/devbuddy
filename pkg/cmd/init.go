package cmd

import (
	"errors"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/ui"
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

	cmdUI := ui.NewTerminal(cfg.DebugEnabled)

	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	return createManifest(cmdUI, projectPath, templateName)
}

func createManifest(cmdUI *ui.UI, projectPath string, templateName string) error {
	templates := manifest.ListTemplates()

	if templateName == "" || !slices.Contains(templates, templateName) {
		options := make([]ui.SelectOption, 0, len(templates))
		for _, template := range templates {
			options = append(options, ui.SelectOption{
				Value: template,
				Label: template,
			})
		}

		result, err := cmdUI.Prompts().Select(ui.SelectRequest{
			Label:   "Select a template",
			Options: options,
		})
		if errors.Is(err, ui.ErrPromptCancelled) {
			return nil
		}
		if err != nil {
			return err
		}
		templateName = result
	}

	if err := manifest.Create(projectPath, templateName); err != nil {
		return err
	}

	cmdUI.ActionHeader("Created dev.yml with template " + templateName)
	cmdUI.ActionNotice("Open dev.yml to adjust for your needs.")
	return nil
}
