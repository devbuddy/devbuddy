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

	// task.addAction(&golangGoPath{})
	task.addActionWithBuilder("", func(ctx *context) error {
		ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
		return nil
	}).addConditionFunc(func(ctx *context) *actionResult {
		if ctx.env.Get("GOPATH") == "" {
			return actionNeeded("GOPATH is not set")
		}
		return actionNotNeeded()
	})

	// task.addAction(&golangInstall{version: version})
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

// type golangGoPath struct{}

// func (g *golangGoPath) description() string {
// 	return ""
// }

// func (g *golangGoPath) needed(ctx *context) *actionResult {
// 	if ctx.env.Get("GOPATH") == "" {
// 		return actionNeeded("GOPATH is not set")
// 	}
// 	return actionNotNeeded()
// }

// func (g *golangGoPath) run(ctx *context) error {
// 	ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
// 	return nil
// }

// type golangInstall struct {
// 	version string
// }

// func (g *golangInstall) description() string {
// 	return fmt.Sprintf("Install Go version %s", g.version)
// }

// func (g *golangInstall) needed(ctx *context) *actionResult {
// 	if !helpers.NewGolang(ctx.cfg, g.version).Exists() {
// 		return actionNeeded("golang distribution is not installed")
// 	}
// 	return actionNotNeeded()
// }

// func (g *golangInstall) run(ctx *context) error {
// 	return helpers.NewGolang(ctx.cfg, g.version).Install()
// }
