package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	api.Register("node", "NodeJS", parseNode)
}

func parseNode(config *api.TaskConfig, task *api.Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.Info = version

	run := func(ctx *context.Context) error {
		return helpers.NewNode(ctx, version).Install()
	}
	condition := func(ctx *context.Context) *api.ActionResult {
		if !helpers.NewNode(ctx, version).Exists() {
			return api.Needed("node version is not installed")
		}
		return api.NotNeeded()
	}
	task.AddActionBuilder("install nodejs from https://nodejs.org", run).
		On(api.FuncCondition(condition)).
		SetFeature("node", version)

	npmInstall := func(ctx *context.Context) error {
		if !utils.PathExists("package.json") {
			ctx.UI.TaskWarning("No package.json found.")
			return nil
		}
		ctx.UI.TaskCommand("npm", "install", "--no-progress")
		return ctx.Executor.Run(executor.New("npm", "install", "--no-progress")).Error
	}
	task.AddActionBuilder("install dependencies with NPM", npmInstall).
		On(api.FileCondition("package.json"))

	return nil
}
