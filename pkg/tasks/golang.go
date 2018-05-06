package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
)

func init() {
	allTasks["go"] = newGolang
}

type Golang struct {
	version string
}

func newGolang(config *taskConfig) (Task, error) {
	task := &Golang{}

	var err error
	task.version, err = config.getPayloadAsString()
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (g *Golang) name() string {
	return "Golang"
}

func (g *Golang) header() string {
	return g.version
}

func (g *Golang) actions(ctx *context) []taskAction {
	return []taskAction{
		&golangGoPath{},
		&golangInstall{version: g.version},
	}
}

func (g *Golang) feature(proj *project.Project) (string, string) {
	return "golang", g.version
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
