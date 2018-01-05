package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Ensure the project is up and running",
	Run:   upRun,
	// Args:  OnlyOneArg,
}

func upRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.NewUI(cfg)

	proj, err := project.FindCurrent()
	checkError(err)

	taskList, err := proj.GetTasks()
	checkError(err)

	for _, task := range taskList {
		err = task.Perform(ui)
		if err != nil {
			break
		}
	}
}
