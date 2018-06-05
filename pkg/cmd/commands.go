package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func customCommandRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui := termui.NewUI(cfg)

	proj, err := project.FindCurrent()
	if err != nil {
		return err
	}

	name := cmd.Annotations["name"]
	spec, ok := proj.Manifest.Commands[name]
	if !ok {
		return fmt.Errorf("custom command is not found: %s", name)
	}

	cmdline := strings.Join(append([]string{spec.Run}, args...), " ")

	ui.CommandHeader(cmdline)

	return executor.NewShell(cmdline).SetCwd(proj.Path).Run()
}

func buildCustomCommands() {
	proj, err := project.FindCurrent()
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
			RunE:               customCommandRun,
			Annotations:        map[string]string{"name": name},
			DisableFlagParsing: true,
			SilenceUsage:       true,
		}
		rootCmd.AddCommand(cmd)
	}
}
