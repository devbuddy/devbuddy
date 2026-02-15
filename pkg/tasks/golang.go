package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("go", "Golang", parseGolang)
}

func parseGolang(config *api.TaskConfig, task *api.Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	if config.IsHash() {
		enabled, present, err := config.GetBooleanProperty("modules")
		if err != nil {
			return err
		}
		if present && !enabled {
			return fmt.Errorf(`"modules: false" is no longer supported for task "go"`)
		}
	}

	task.Info = version

	installNeeded := func(ctx *context.Context) *api.ActionResult {
		if !helpers.NewGolang(ctx, version).Exists() {
			return api.Needed("golang distribution is not installed")
		}
		return api.NotNeeded()
	}
	installGo := func(ctx *context.Context) error {
		return helpers.NewGolang(ctx, version).Install()
	}
	task.AddActionBuilder("install golang distribution", installGo).
		On(api.FuncCondition(installNeeded)).
		SetFeature("golang", version)

	return nil
}
