package tasks

import (
	"fmt"
	"path"
	"strings"

	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/helpers"
)

func init() {
	allTasks["pip"] = NewPip
}

type Pip struct {
	files []string
}

func NewPip() Task {
	return &Pip{}
}

func (p *Pip) Load(definition interface{}) (bool, error) {
	def, ok := definition.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}

	if payload, ok := def["pip"]; ok {
		for _, value := range payload.([]interface{}) {
			if v, ok := value.(string); ok {
				p.files = append(p.files, v)

			} else {
				return false, fmt.Errorf("invalid pip files")
			}
		}
		if len(p.files) > 0 {
			return true, nil
		}

		return false, fmt.Errorf("no pip files specified")
	}
	return false, nil
}

func (p *Pip) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader("Pip", strings.Join(p.files, ", "))

	pipCmd := "pip"

	pythonParam, hasPythonFeature := ctx.features["python"]
	if hasPythonFeature {
		venv := helpers.NewVirtualenv(ctx.cfg, pythonParam)
		pipCmd = path.Join(venv.BinPath(), "pip")
	}

	for _, file := range p.files {
		code, err := executor.Run(pipCmd, "install", "--require-virtualenv", "-r", file)
		if err != nil {
			return err
		}
		if code != 0 {
			return fmt.Errorf("Pip failed with code %s", code)
		}
	}
	return nil
}
