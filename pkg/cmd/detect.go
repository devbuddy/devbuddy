package cmd

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect project tasks",
	Run:   detectRun,
	Args:  noArgs,
}

func detectRun(cmd *cobra.Command, args []string) {
	proj, err := project.FindCurrent()
	checkError(err)

	projectTasks, err := tasks.DetectTasksFromProject(proj)
	checkError(err)

	fmt.Print(tasks.InspectTasks(projectTasks, proj))
}
