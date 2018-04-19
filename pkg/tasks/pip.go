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

func newPip(config *taskConfig) (Task, error) {
	task := &Pip{}

	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			task.files = append(task.files, v)
		} else {
			return nil, fmt.Errorf("invalid pip files")
		}
	}
	if len(task.files) == 0 {
		return nil, fmt.Errorf("no pip files specified")
	}

	return task, nil
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
