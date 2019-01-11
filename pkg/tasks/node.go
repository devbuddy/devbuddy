package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.RegisterTaskDefinition("node", "NodeJS", parseNode)
}

func parseNode(config *taskapi.TaskConfig, task *taskapi.Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.Header = version
	task.feature = features.NewFeatureInfo("node", version)

	run := func(ctx *taskapi.Context) error {
		return helpers.NewNode(ctx.cfg, version).Install()
	}
	condition := func(ctx *taskapi.Context) *taskapi.ActionResult {
		if !helpers.NewNode(ctx.cfg, version).Exists() {
			return actionNeeded("node version is not installed")
		}
		return actionNotNeeded()
	}
	taskapi.AddActionWithBuilder("install nodejs from https://nodejs.org", run).OnFunc(condition)

	return nil
}
