package hook

import (
	"fmt"
	"os"
	"time"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/features"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
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

	ui.Debug("total time: %s", time.Now().Sub(timerStart))
}

func handleFeatures(cfg *config.Config, proj *project.Project, ui *termui.HookUI) {
	env := features.NewEnv(os.Environ())

	runner := features.NewRunner(cfg, proj, ui, env)
	runner.Run()

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
