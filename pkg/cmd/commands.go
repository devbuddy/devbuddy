package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/project"
)

func findProject() (*project.Project, error) {
	path, err := os.Getwd()
	checkError(err)

	return project.FindCurrent(path)
}

func customCommandRun(cmd *cobra.Command, args []string) {
	proj, err := findProject()
	checkError(err)

	name := cmd.Annotations["name"]
	spec, ok := proj.Manifest.Commands[name]
	if !ok {
		exitWithMessage(fmt.Sprintf("custom command is not found: %s", name))
	}

	err = executor.RunShell(spec.Run)
	if err != nil {
		fmt.Printf("Command failed: %s", err)
	}
}

func buildCustomCommands() {
	proj, err := findProject()
	if err != nil {
		return
	}

	var cmd *cobra.Command

	for name, spec := range proj.Manifest.Commands {
		desc := "Custom"
		if spec.Description != "" {
			desc = fmt.Sprintf("Custom: %s", spec.Description)
		}

		useLine := fmt.Sprintf("%s [ARGS...]", name)

		cmd = &cobra.Command{
			Use:                useLine,
			Short:              desc,
			Run:                customCommandRun,
			Annotations:        map[string]string{"name": name},
			DisableFlagParsing: true,
		}
		rootCmd.AddCommand(cmd)
	}
}
