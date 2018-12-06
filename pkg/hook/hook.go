package hook

import (
	"fmt"
	"time"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func Hook() {
	// In the shell hook, the stdout is evaluated by the shell
	// stderr is used to display messages to the user

	// Also, we can't annoy the user here, so we always just quit silently

	timerStart := time.Now()

	cfg, configErr := config.Load()
	ui := termui.NewHookUI(cfg)

	if configErr != nil {
		ui.Debug("error while loading the config: %s", configErr)
		return
	}

	proj, err := project.FindCurrent()
	if err != nil && err != project.ErrProjectNotFound {
		ui.Debug("error while searching for current project: %s", err)
		return
	}

	handleFeatures(cfg, proj, ui)

	ui.Debug("total time: %s", time.Since(timerStart))
}

func handleFeatures(cfg *config.Config, proj *project.Project, ui *termui.HookUI) {
	allFeatures, err := getFeaturesFromProject(proj)
	if err != nil {
		ui.Debug("error while building the project tasks: %s", err)
		return
	}

	env := env.NewFromOS()

	runner := features.NewRunner(cfg, proj, ui, env)
	runner.Run(allFeatures)

	envChanges := env.Changed()
	for _, change := range envChanges {
		ui.Debug("Env change: %+v", change)

		if change.Deleted {
			fmt.Printf("unset %s\n", change.Name)
		} else {
			fmt.Printf("export %s=\"%s\"\n", change.Name, change.Value)
		}
	}
}

func getFeaturesFromProject(proj *project.Project) (features map[string]string, err error) {
	if proj == nil {
		return map[string]string{}, nil
	}

	var taskList []*tasks.Task

	if manifest.ExistsIn(proj.Path) {
		taskList, err = tasks.GetTasksFromProjectManifest(proj)
	}

	if err != nil {
		return nil, err
	}

	return tasks.GetFeaturesFromTasks(proj, taskList), nil
}
