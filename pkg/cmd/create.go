package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/integration"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

var createCmd = &cobra.Command{
	Use:   "create [PROJECT]",
	Short: "Create a new project",
	Run:   createRun,
	Args:  onlyOneArg,
}

func createRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.NewUI(cfg)

	proj, err := project.NewFromID(args[0], cfg)
	checkError(err)

	if proj.Exists() {
		ui.ProjectExists()
	} else {
		err = proj.Create()
		checkError(err)
	}

	ui.JumpProject(proj.FullName())
	err = integration.AddFinalizerCd(proj.Path)
	checkError(err)
}
