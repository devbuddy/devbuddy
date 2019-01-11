package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	t := registerTaskDefinition("go")
	t.name = "Golang"
	t.parser = parseGolang
}

func parseGolang(config *TaskConfig, task *Task) error {
	version, err := config.getStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.header = version
	task.feature = features.NewFeatureInfo("golang", version)

	checkPATHVar := func(ctx *Context) *actionResult {
		if ctx.env.Get("GOPATH") == "" {
			return actionNeeded("GOPATH is not set")
		}
		return actionNotNeeded()
	}
	showPATHWarning := func(ctx *Context) error {
		ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	}
	task.AddActionWithBuilder("", showPATHWarning).OnFunc(checkPATHVar)

	installNeeded := func(ctx *Context) *actionResult {
		if !helpers.NewGolang(ctx.cfg, version).Exists() {
			return actionNeeded("golang distribution is not installed")
		}
		return actionNotNeeded()
	}
	installGo := func(ctx *Context) error {
		return helpers.NewGolang(ctx.cfg, version).Install()
	}
	task.AddActionWithBuilder("install golang distribution", installGo).OnFunc(installNeeded)

	return nil
}
