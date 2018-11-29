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

	action := newAction("", func(ctx *context) error {
		ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	})
	action.onFunc(func(ctx *context) *actionResult {
		if ctx.env.Get("GOPATH") == "" {
			return actionNeeded("GOPATH is not set")
		}
		return actionNotNeeded()
	})
	action.onFileChange("setup.py")
	action.onFeatureChange("feature")
	task.addAction(action)

	action = newAction("", func(ctx *context) error {
		ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	}).onFunc(func(ctx *context) *actionResult {
		if ctx.env.Get("GOPATH") == "" {
			return actionNeeded("GOPATH is not set")
		}
		return actionNotNeeded()
	})
	task.addAction(action)

	action = newAction(fmt.Sprintf("Install Go version %s", version),
		func(ctx *context) error {
			return helpers.NewGolang(ctx.cfg, version).Install()
		}).
		onFunc(func(ctx *context) *actionResult {
			if !helpers.NewGolang(ctx.cfg, version).Exists() {
				return actionNeeded("golang distribution is not installed")
			}
			return actionNotNeeded()
		})
	task.addAction(action)

	return nil
}
