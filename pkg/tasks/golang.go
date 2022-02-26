package tasks

import (
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

	modulesEnabled := false
	if config.IsHash() {
		modulesEnabled, err = config.GetBooleanPropertyDefault("modules", false)
		if err != nil {
			return err
		}
	}
	featureVersion := version
	if modulesEnabled {
		featureVersion += "+mod"
	}

	task.Info = version

	checkPATHVar := func(ctx *context.Context) *api.ActionResult {
		if ctx.Env.Get("GOPATH") == "" {
			return api.Needed("GOPATH is not set")
		}
		return api.NotNeeded()
	}
	showPATHWarning := func(ctx *context.Context) error {
		ctx.UI.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	}
	task.AddActionBuilder("", showPATHWarning).On(api.FuncCondition(checkPATHVar))

	installNeeded := func(ctx *context.Context) *api.ActionResult {
		if !helpers.NewGolang(ctx.Cfg, version).Exists() {
			return api.Needed("golang distribution is not installed")
		}
		return api.NotNeeded()
	}
	installGo := func(ctx *context.Context) error {
		return helpers.NewGolang(ctx.Cfg, version).Install()
	}
	task.AddActionBuilder("install golang distribution", installGo).
		On(api.FuncCondition(installNeeded)).
		SetFeature("golang", featureVersion)

	return nil
}
