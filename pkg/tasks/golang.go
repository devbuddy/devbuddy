package tasks

import (
	"os"

	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
)

func init() {
	allTasks["go"] = NewGolang
}

type Golang struct {
	version string
}

func NewGolang() Task {
	return &Golang{}
}

func (g *Golang) Load(config *taskConfig) (bool, error) {
	version, ok := config.payload.(string)
	if !ok {
		return false, nil
	}
	g.version = version
	return true, nil
}

func (g *Golang) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader("Golang", g.version)

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
