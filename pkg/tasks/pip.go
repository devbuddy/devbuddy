package tasks

import (
	"fmt"
	"strings"

	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/helpers"
)

func init() {
	allTasks["pip"] = newPip
}

type Pip struct {
	files []string
}

func newPip() Task {
	return &Pip{}
}

func (p *Pip) load(config *taskConfig) error {
	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			p.files = append(p.files, v)

		} else {
			return fmt.Errorf("invalid pip files")
		}
	}
	if len(p.files) > 0 {
		return nil
	}

	return fmt.Errorf("no pip files specified")
}

func (p *Pip) name() string {
	return "Pip"
}

func (p *Pip) header() string {
	return strings.Join(p.files, ", ")
}

func (p *Pip) perform(ctx *Context) (err error) {
	// We should also check that the python task is executed before this one
	pythonParam, hasPythonFeature := ctx.features["python"]
	if !hasPythonFeature {
		return fmt.Errorf("You must specify a Python environment to use this task")
	}

	pythonCmd := helpers.NewVirtualenv(ctx.cfg, pythonParam).Which("python")

	for _, file := range p.files {
		code, err := executor.Run(pythonCmd, "-m", "pip", "install", "-r", file)
		if err != nil {
			return err
		}
		if code != 0 {
			return fmt.Errorf("Pip failed with code %d", code)
		}
	}
	return nil
}
