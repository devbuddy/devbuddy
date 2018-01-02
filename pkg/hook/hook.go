package hook

import (
	"fmt"
	"os"

	"github.com/pior/dad/pkg/features"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

func Hook() {
	// In the shell hook, the stdout is evaluated by the shell
	// stderr is used to display messages to the user
	ui := termui.NewHookUI()

	proj, err := project.FindCurrent()

	if err != nil && err != project.ErrProjectNotFound {
		// We can't annoy the user here, just quit silently
		ui.Debug("error while searching for current project: %s", err)
		return
	}

	handleFeatures(proj, ui)
}

func handleFeatures(proj *project.Project, ui *termui.HookUI) {
	env := features.NewEnv(os.Environ())

	runner := features.NewRunner(proj, ui, env)
	runner.Run()

	envChanges := env.Changed()
	for _, change := range envChanges {
		ui.Debug("Env change: %+v", change)

		if change.Deleted {
			fmt.Printf("unset %s", change.Name)
		} else {
			fmt.Printf("export %s=\"%s\"", change.Name, change.Value)
		}
	}
}
