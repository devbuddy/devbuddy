package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Ensure the project is up and running",
	Run:   upRun,
	Args:  noArgs,
}

func upRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.New(cfg)

	proj, err := project.FindCurrent()
	checkError(err)

	taskList, err := tasks.GetTasksFromProject(proj)
	checkError(err)

	ctx := tasks.NewContext(cfg, proj, ui, taskList)

	success, err := tasks.Run(ctx, &tasks.TaskRunnerImpl{}, taskList)
	checkError(err)
	if !success {
		os.Exit(1)
	}
}
