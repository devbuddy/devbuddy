package tasks

import (
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

func (g *Golang) Load(definition interface{}) (bool, error) {
	def, ok := definition.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}
	if version, ok := def["go"]; ok {
		g.version, ok = version.(string)
		if !ok {
			return false, nil
		}
		return true, nil
	}

	return false, nil
}

func (g *Golang) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader("Golang", g.version)

	goSrc := helpers.NewGolang(ctx.cfg, g.version)

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
