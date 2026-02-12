package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/integration"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

var cdCmd = &cobra.Command{
	Use:          "cd [PROJECT]",
	Short:        "Jump to a local project",
	RunE:         cdRun,
	Args:         onlyOneArg,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func cdRun(_ *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui := termui.New(cfg)

	proj, err := project.FindBestMatch(args[0], cfg)
	if err != nil {
		return err
	}

	ui.JumpProject(proj.FullName())

	return integration.AddFinalizerCd(proj.Path)
}
