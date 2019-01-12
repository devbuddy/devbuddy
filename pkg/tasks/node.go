package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	Register("node", "NodeJS", parseNode)
}

func parseNode(config *TaskConfig, task *Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.SetInfo(version)
	task.SetFeature("node", version)

	run := func(ctx *Context) error {
		return helpers.NewNode(ctx.Cfg, version).Install()
	}
	condition := func(ctx *Context) *ActionResult {
		if !helpers.NewNode(ctx.Cfg, version).Exists() {
			return ActionNeeded("node version is not installed")
		}
		return ActionNotNeeded()
	}
	task.AddActionWithBuilder("install nodejs from https://nodejs.org", run).OnFunc(condition)

	return nil
}
