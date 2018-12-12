package tasks

import "github.com/devbuddy/devbuddy/pkg/helpers"

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

	return nil
}
