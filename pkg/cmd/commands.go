package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func customCommandRun(cmd *cobra.Command, args []string) error {
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	man, err := manifest.Load(ctx.Project.Path)
	if err != nil {
		return err
	}

	cmds, err := man.GetCommands()
	if err != nil {
		return err
	}

	name := cmd.Annotations["name"]
	spec, ok := cmds[name]
	if !ok {
		return fmt.Errorf("custom command is not found: %s", name)
	}

	cmdline := strings.Join(append([]string{spec.Run}, args...), " ")

	ctx.UI.CommandHeader(cmdline)

	for name, value := range man.Env {
		if !ctx.Env.Has(name) {
			ctx.Env.Set(name, value)
		}
	}

	execCmd := executor.NewShell(cmdline)
	execCmd.Passthrough = true
	return ctx.Executor.Run(execCmd).Error
}

func buildCustomCommands(rootCmd *cobra.Command) {
	proj, err := project.FindCurrent()
	if err != nil {
		return
	}

	man, err := manifest.Load(proj.Path)
	if err != nil {
		return
	}

	cmds, err := man.GetCommands()
	if err != nil {
		return
	}

	for name, spec := range cmds {
		rootCmd.AddCommand(&cobra.Command{
			Use:                fmt.Sprintf("%s [ARGS...]", name),
			Short:              spec.Description,
			RunE:               customCommandRun,
			Annotations:        map[string]string{"name": name},
			DisableFlagParsing: true,
			SilenceUsage:       true,
			GroupID:            "project",
		})
	}
}
