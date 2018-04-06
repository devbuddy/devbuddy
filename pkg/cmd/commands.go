package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
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

	code, err := executor.NewShell(cmdline).SetCwd(proj.Path).Run()
	if err != nil {
		return fmt.Errorf("command failed: %s", err)
	}
	if code != 0 {
		return fmt.Errorf("command failed with code %d", code)
	}

	return nil
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
