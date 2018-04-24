package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/integration"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

var cloneCmd = &cobra.Command{
	Use:   "clone [REMOTE]",
	Short: "Clone a project from github.com",
	Run:   cloneRun,
	Args:  onlyOneArg,
}

func cloneRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.NewUI(cfg)

	proj, err := project.NewFromID(args[0], cfg)
	checkError(err)

	if proj.Exists() {
		ui.ProjectExists()
	} else {
		err = proj.Clone()
		checkError(err)
	}

	ui.JumpProject(proj.FullName())
	err = integration.AddFinalizerCd(proj.Path)
	checkError(err)
}
