package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.RegisterTaskDefinition("go", "Golang", parseGolang)
}

func parseGolang(config *taskapi.TaskConfig, task *taskapi.Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.Header = version
	task.Feature = features.NewFeatureInfo("golang", version)

	checkPATHVar := func(ctx *taskapi.Context) *taskapi.ActionResult {
		if ctx.env.Get("GOPATH") == "" {
			return taskapi.ActionNeeded("GOPATH is not set")
		}
		return actionNotNeeded()
	}
	showPATHWarning := func(ctx *taskapi.Context) error {
		ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	}
	task.AddActionWithBuilder("", showPATHWarning).OnFunc(checkPATHVar)

	installNeeded := func(ctx *taskapi.Context) *taskapi.ActionResult {
		if !helpers.NewGolang(ctx.cfg, version).Exists() {
			return taskapi.ActionNeeded("golang distribution is not installed")
		}
		return actionNotNeeded()
	}
	installGo := func(ctx *taskapi.Context) error {
		return helpers.NewGolang(ctx.cfg, version).Install()
	}
	task.AddActionWithBuilder("install golang distribution", installGo).OnFunc(installNeeded)

	return nil
}
