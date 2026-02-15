package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func customCommandRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui := termui.New(cfg)

	proj, err := project.FindCurrent()
	if err != nil {
		return err
	}

	man, err := manifest.Load(proj.Path)
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

	ui.CommandHeader(cmdline)

	exec := executor.NewShell(cmdline).SetPassthrough(true).SetCwd(proj.Path)

	envs := env.NewFromOS()
	for name, value := range man.Env {
		if !envs.Has(name) {
			envs.Set(name, value)
		}
	}
	exec.SetEnv(envs.Environ())

	return exec.Run().Error
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
