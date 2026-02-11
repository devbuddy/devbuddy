package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

var inspectCmd = &cobra.Command{
	Use:          "inspect",
	Short:        "Inspect the project and its tasks",
	RunE:         inspectRun,
	Args:         noArgs,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func inspectRun(_ *cobra.Command, _ []string) error {
	proj, err := project.FindCurrent()
	if err != nil {
		return err
	}

	fmt.Printf("Found project at %s\n", proj.Path)

	projectTasks, err := api.GetTasksFromProject(proj)
	if err != nil {
		return err
	}

	for _, task := range projectTasks {
		fmt.Printf("- %s\n", task.Describe())
	}
	return nil
}
