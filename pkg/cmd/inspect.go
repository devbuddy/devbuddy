package cmd

import (
	"fmt"

	"github.com/pior/dad/pkg/tasks"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/project"
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
	fmt.Printf("Manifest: %s\n", proj.Manifest.Path)

	projectTasks, err := tasks.GetTasksFromProject(proj)
	checkError(err)

	fmt.Print(tasks.InspectTasks(projectTasks, proj))
}
