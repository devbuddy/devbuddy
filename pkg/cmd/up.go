package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskengine"
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
	ctx, err := context.Load()
	checkError(err)

	taskList, err := taskapi.GetTasksFromProject(ctx.Project)
	checkError(err)

	runner := &taskengine.TaskRunnerImpl{}
	selector := taskengine.NewTaskSelector()

	success, err := taskengine.Run(ctx, runner, selector, taskList)
	checkError(err)
	if !success {
		os.Exit(1)
	}
}
