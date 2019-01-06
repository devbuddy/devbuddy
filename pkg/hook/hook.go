package hook

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func Run() {
	// In the shell hook, the stdout is evaluated by the shell
	// stderr is used to display messages to the user

	// Also, we can't annoy the user here, so we always just quit silently

	cfg, err := config.Load()
	if err != nil {
		return
	}

	ui := termui.NewHook(cfg)

	err = run(cfg, ui)
	if err != nil {
		ui.Debug("%s", err)
	}
}

func run(cfg *config.Config, ui *termui.UI) error {
	proj, err := project.FindCurrent()
	if err != nil {
		return err
	}

	allTasks, err := tasks.GetTasksFromProject(proj)
	if err != nil {
		return err
	}

	env := env.NewFromOS()
	features.Sync(cfg, proj, ui, env, tasks.GetFeaturesFromTasks(allTasks))
	printEnvironmentChangeAsShellCommands(ui, env)

	return nil
}

func printEnvironmentChangeAsShellCommands(ui *termui.UI, env *env.Env) {
	for _, change := range env.Changed() {
		ui.Debug("Env change: %+v", change)

		if change.Deleted {
			fmt.Printf("unset %s\n", change.Name)
		} else {
			fmt.Printf("export %s=\"%s\"\n", change.Name, change.Value)
		}
	}
}
