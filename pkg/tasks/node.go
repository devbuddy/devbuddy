package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("node", "NodeJS", parseNode)
}

func parseNode(config *taskapi.TaskConfig, task *taskapi.Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.Info = version
	task.SetFeature("node", version)

	run := func(ctx *taskapi.Context) error {
		return helpers.NewNode(ctx.Cfg, version).Install()
	}
	condition := func(ctx *taskapi.Context) *taskapi.ActionResult {
		if !helpers.NewNode(ctx.Cfg, version).Exists() {
			return taskapi.ActionNeeded("node version is not installed")
		}
		return taskapi.ActionNotNeeded()
	}
	task.AddActionWithBuilder("install nodejs from https://nodejs.org", run).OnFunc(condition)

	return nil
}
