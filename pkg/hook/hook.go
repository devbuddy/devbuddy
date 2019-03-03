package hook

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func Run() {
	// In the shell hook, the stdout is evaluated by the shell
	// stderr is used to display messages to the user

	// Also, we can't annoy the user here, so we always just quit silently

	ctx, err := context.Load()
	if err != nil {
		return
	}
	ctx.UI.SetOutputToStderr()

	err = run(ctx)
	if err != nil {
		ctx.UI.Debug("%s", err)
	}
}

func run(ctx *context.Context) error {
	allFeatures, err := getFeaturesFromProject(ctx.Project)
	if err != nil {
		return err
	}
	ctx.UI.Debug("features: %+v", allFeatures)

	autoenv.Sync(ctx, allFeatures)
	printEnvironmentChangeAsShellCommands(ctx.UI, ctx.Env)

	return nil
}

func getFeaturesFromProject(proj *project.Project) (autoenv.FeatureSet, error) {
	if proj == nil {
		// When no project was found, we must deactivate all potentially active features
		// So we continue with an empty feature map
		return autoenv.NewFeatureSet(), nil
	}
	allTasks, err := taskapi.GetTasksFromProject(proj)
	if err != nil {
		return nil, err
	}
	return taskapi.GetFeaturesFromTasks(allTasks), nil
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
