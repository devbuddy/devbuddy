package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/integration"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

var createCmd = &cobra.Command{
	Use:     "create [PROJECT]",
	Short:   "Create a new project",
	Run:     createRun,
	Args:    onlyOneArg,
	GroupID: "devbuddy",
}

func createRun(_ *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.New(cfg)

	proj, err := project.NewFromID(args[0], cfg)
	checkError(err)

	if proj.Exists() {
		ui.ProjectExists()
	} else {
		err = proj.Create()
		checkError(err)

		err = createManifest(ui, proj.Path)
		checkError(err)
	}

	ui.JumpProject(proj.FullName())
	err = integration.AddFinalizerCd(proj.Path)
	checkError(err)
}
