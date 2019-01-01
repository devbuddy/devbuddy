package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	t := registerTaskDefinition("node")
	t.name = "NodeJS"
	t.parser = parseNode
}

func parseNode(config *taskConfig, task *Task) error {
	version, err := config.getStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.header = version
	task.featureName = "node"
	task.featureParam = version

	builder := actionBuilder("install nodejs from https://nodejs.org", func(ctx *Context) error {
		return helpers.NewNode(ctx.cfg, version).Install()
	})
	builder.OnFunc(func(ctx *Context) *actionResult {
		if !helpers.NewNode(ctx.cfg, version).Exists() {
			return actionNeeded("node version is not installed")
		}
		return actionNotNeeded()
	})
	task.addAction(builder.Build())

	builder = actionBuilder("install dependencies", func(ctx *Context) error {
		if !utils.PathExists("package.json") {
			ctx.ui.TaskWarning("No package.json found.")
			return nil
		}
		return command(ctx, "npm", "install", "--no-progress").Run().Error
	})
	builder.OnFileChange("package.json")
	task.addAction(builder.Build())

	return nil
}
