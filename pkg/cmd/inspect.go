package cmd

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/tasks"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/project"
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

	projectTasks, err := tasks.GetTasksFromProject(proj)
	checkError(err)

	fmt.Print(tasks.InspectTasks(projectTasks, proj))
}
