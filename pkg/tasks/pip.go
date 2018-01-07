package tasks

import (
	"fmt"
	"strings"

	"github.com/pior/dad/pkg/executor"
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

	for _, file := range p.files {
		_, err = executor.RunSilent("pip", "install", "-r", file)
		if err != nil {
			return
		}
	}
	return nil
}
