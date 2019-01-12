package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	Register("go", "Golang", parseGolang)
}

func parseGolang(config *TaskConfig, task *Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.SetInfo(version)
	task.SetFeature("golang", version)

	checkPATHVar := func(ctx *Context) *ActionResult {
		if ctx.env.Get("GOPATH") == "" {
			return ActionNeeded("GOPATH is not set")
		}
		return ActionNotNeeded()
	}
	showPATHWarning := func(ctx *Context) error {
		ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	}
	task.AddActionWithBuilder("", showPATHWarning).OnFunc(checkPATHVar)

	installNeeded := func(ctx *Context) *ActionResult {
		if !helpers.NewGolang(ctx.cfg, version).Exists() {
			return ActionNeeded("golang distribution is not installed")
		}
		return ActionNotNeeded()
	}
	installGo := func(ctx *Context) error {
		return helpers.NewGolang(ctx.cfg, version).Install()
	}
	task.AddActionWithBuilder("install golang distribution", installGo).OnFunc(installNeeded)

	return nil
}
