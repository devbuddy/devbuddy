package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskengine"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func init() {
	tasks.RegisterTasks()
}

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

	taskList, err := taskapi.GetTasksFromProject(proj)
	checkError(err)

	ctx := &taskapi.Context{
		Cfg:      cfg,
		Project:  proj,
		UI:       ui,
		Env:      env.NewFromOS(),
		Features: taskapi.GetFeaturesFromTasks(taskList),
	}

	runner := &taskengine.TaskRunnerImpl{}
	selector := taskengine.NewTaskSelector()

	success, err := taskengine.Run(ctx, runner, selector, taskList)
	checkError(err)
	if !success {
		os.Exit(1)
	}
}
