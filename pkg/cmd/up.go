package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskengine"
)

func init() {
	tasks.RegisterTasks()
}

var upCmd = &cobra.Command{
	Use:          "up",
	Short:        "Ensure the project is up and running",
	RunE:         upRun,
	Args:         noArgs,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func upRun(_ *cobra.Command, _ []string) error {
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	taskList, err := api.GetTasksFromProject(ctx.Project)
	if err != nil {
		return err
	}

	runner := taskengine.NewTaskRunner(ctx)
	selector := taskengine.NewTaskSelector()

	success, err := taskengine.Run(ctx, runner, selector, taskList)
	if err != nil {
		return err
	}
	if !success {
		return errTasksFailed
	}
	return nil
}
