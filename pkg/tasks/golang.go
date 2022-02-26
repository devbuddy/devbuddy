package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("go", "Golang", parseGolang)
}

func parseGolang(config *taskapi.TaskConfig, task *taskapi.Task) error {
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

	checkPATHVar := func(ctx *context.Context) *taskapi.ActionResult {
		if ctx.Env.Get("GOPATH") == "" {
			return taskapi.ActionNeeded("GOPATH is not set")
		}
		return taskapi.ActionNotNeeded()
	}
	showPATHWarning := func(ctx *context.Context) error {
		ctx.UI.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	}
	task.AddActionBuilder("", showPATHWarning).On(taskapi.FuncCondition(checkPATHVar))

	installNeeded := func(ctx *context.Context) *taskapi.ActionResult {
		if !helpers.NewGolang(ctx.Cfg, version).Exists() {
			return taskapi.ActionNeeded("golang distribution is not installed")
		}
		return taskapi.ActionNotNeeded()
	}
	installGo := func(ctx *context.Context) error {
		return helpers.NewGolang(ctx.Cfg, version).Install()
	}
	task.AddActionBuilder("install golang distribution", installGo).
		On(taskapi.FuncCondition(installNeeded)).
		SetFeature("golang", featureVersion)

	return nil
}
