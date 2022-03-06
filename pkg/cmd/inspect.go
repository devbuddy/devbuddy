package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect the project and its tasks",
	Run:   inspectRun,
	Args:  noArgs,
}

func inspectRun(cmd *cobra.Command, args []string) {
	proj, err := project.FindCurrent()
	checkError(err)

	fmt.Printf("Found project at %s\n", proj.Path)

	projectTasks, err := api.GetTasksFromProject(proj)
	checkError(err)

	for _, task := range projectTasks {
		fmt.Printf("- %s\n", task.Describe())
	}
}
