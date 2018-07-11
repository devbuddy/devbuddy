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

	task.addAction(&golangGoPath{})
	task.addAction(&golangInstall{version: version})

	return nil
}

type golangGoPath struct{}

func (g *golangGoPath) description() string {
	return ""
}

func (g *golangGoPath) needed(ctx *context) (bool, error) {
	return ctx.env.Get("GOPATH") == "", nil
}

func (g *golangGoPath) run(ctx *context) error {
	ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
	return nil
}

type golangInstall struct {
	version string
}

func (g *golangInstall) description() string {
	return fmt.Sprintf("Install Go version %s", g.version)
}

func (g *golangInstall) needed(ctx *context) (bool, error) {
	return !helpers.NewGolang(ctx.cfg, g.version).Exists(), nil
}

func (g *golangInstall) run(ctx *context) error {
	return helpers.NewGolang(ctx.cfg, g.version).Install()
}
