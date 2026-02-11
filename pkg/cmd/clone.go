package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/integration"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

var cloneCmd = &cobra.Command{
	Use:          "clone [REMOTE]",
	Short:        "Clone a project from github.com",
	RunE:         cloneRun,
	Args:         onlyOneArg,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func cloneRun(_ *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui := termui.New(cfg)

	proj, err := project.NewFromID(args[0], cfg)
	if err != nil {
		return err
	}

	if proj.Exists() {
		ui.ProjectExists()
	} else {
		if err := proj.Clone(); err != nil {
			return err
		}
	}

	ui.JumpProject(proj.FullName())
	return integration.AddFinalizerCd(proj.Path)
}
