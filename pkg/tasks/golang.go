package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	t := registerTaskDefinition("go")
	t.name = "Golang"
	t.parser = parseGolang
}

func parseGolang(config *taskConfig, task *Task) error {
	version, err := config.getStringProperty("version", true)
	if err != nil {
		return err
	}

	task.header = version
	task.featureName = "golang"
	task.featureParam = version

	task.addActionWithBuilder("", func(ctx *context) error {
		ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	}).addConditionFunc(func(ctx *context) *actionResult {
		if ctx.env.Get("GOPATH") == "" {
			return actionNeeded("GOPATH is not set")
		}
		return actionNotNeeded()
	})

	task.addActionWithBuilder(fmt.Sprintf("Install Go version %s", version),
		func(ctx *context) error {
			return helpers.NewGolang(ctx.cfg, version).Install()
		}).
		addConditionFunc(func(ctx *context) *actionResult {
			if !helpers.NewGolang(ctx.cfg, version).Exists() {
				return actionNeeded("golang distribution is not installed")
			}
			return actionNotNeeded()
		})

	return nil
}
