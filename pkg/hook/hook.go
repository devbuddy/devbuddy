package hook

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func Run() {
	// In the shell hook, the stdout is evaluated by the shell
	// stderr is used to display messages to the user

	// Also, we can't annoy the user here, so we always just quit silently

	ctx, err := context.Load(true)
	if err != nil {
		return
	}

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
	ctx.UI.Debug("Desired features: %+v", allFeatures)

	autoenv.Sync(ctx, allFeatures)
	emitEnvironmentChangeAsShellCommands(ctx)
	emitShellHashResetCommand(ctx)

	return nil
}

func getFeaturesFromProject(proj *project.Project) (autoenv.FeatureSet, error) {
	if proj == nil {
		// When no project was found, we must deactivate all potentially active features
		// So we continue with an empty feature map
		return autoenv.NewFeatureSet(), nil
	}
	allTasks, err := api.GetTasksFromProject(proj)
	if err != nil {
		return nil, err
	}
	return api.GetFeaturesFromTasks(allTasks), nil
}

func emitEnvironmentChangeAsShellCommands(ctx *context.Context) {
	for _, mutation := range ctx.Env.Mutations() {
		ctx.UI.Debug("Apply change: %s\n%s", mutation.Name, mutation.DiffString())

		if mutation.Current == nil {
			fmt.Printf("unset %s\n", mutation.Name)
		} else {
			escaped := strings.ReplaceAll(mutation.Current.Value, "\"", "\\\"")
			fmt.Printf("export %s=\"%s\"\n", mutation.Name, escaped)
		}
	}
}

func emitShellHashResetCommand(ctx *context.Context) {
	for _, mutation := range ctx.Env.Mutations() {
		if mutation.Name == "PATH" {
			fmt.Printf("hash -r")
			return
		}
	}
}
