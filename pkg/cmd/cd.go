package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/integration"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

var cdCmd = &cobra.Command{
	Use:   "cd [PROJECT]",
	Short: "Jump to a local project",
	Run:   cdRun,
	Args:  onlyOneArg,
}

func cdRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.New(cfg)

	proj, err := project.FindBestMatch(args[0], cfg)
	checkError(err)

	ui.JumpProject(proj.FullName())

	err = integration.AddFinalizerCd(proj.Path)
	checkError(err)
}
