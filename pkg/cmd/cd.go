package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/integration"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
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

	ui := termui.NewUI(cfg)

	proj, err := project.FindBestMatch(args[0], cfg)
	checkError(err)

	ui.JumpProject(proj.FullName())

	err = integration.AddFinalizerCd(proj.Path)
	checkError(err)
}
