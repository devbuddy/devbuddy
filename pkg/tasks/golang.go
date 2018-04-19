package tasks

import (
	"os"

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

func (g *Golang) perform(ctx *Context) (err error) {
	goSrc := helpers.NewGolang(ctx.cfg, g.version)

	if os.Getenv("GOPATH") == "" {
		ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
	}

	if goSrc.Exists() {
		ctx.ui.TaskAlreadyOk()
		return nil
	}

	err = goSrc.Install()
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	ctx.ui.TaskActed()
	return nil
}

func (g *Golang) Feature(proj *project.Project) (string, string) {
	return "golang", g.version
}
