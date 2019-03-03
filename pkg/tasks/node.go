package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/utils"
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

	run := func(ctx *context.Context) error {
		return helpers.NewNode(ctx.Cfg, version).Install()
	}
	condition := func(ctx *context.Context) *taskapi.ActionResult {
		if !helpers.NewNode(ctx.Cfg, version).Exists() {
			return taskapi.ActionNeeded("node version is not installed")
		}
		return taskapi.ActionNotNeeded()
	}
	task.AddActionWithBuilder("install nodejs from https://nodejs.org", run).
		OnFunc(condition).
		SetFeature("node", version)

	npmInstall := func(ctx *context.Context) error {
		if !utils.PathExists("package.json") {
			ctx.UI.TaskWarning("No package.json found.")
			return nil
		}
		return command(ctx, "npm", "install", "--no-progress").Run().Error
	}
	task.AddActionWithBuilder("install dependencies with NPM", npmInstall).
		OnFileChange("package.json")

	return nil
}
