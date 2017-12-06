package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/executor"
)

type Pip struct {
	files []string
}

func (p *Pip) Load(definition map[interface{}]interface{}) (bool, error) {
	if payload, ok := definition["pip"]; ok {
		for _, value := range payload.([]interface{}) {
			if v, ok := value.(string); ok {
				p.files = append(p.files, v)

			} else {
				return false, fmt.Errorf("invalid pip files")
			}
		}
		if len(p.files) > 0 {
			return true, nil
		} else {
			return false, fmt.Errorf("no pip files specified")
		}
	}
	return false, nil
}

func (p *Pip) Perform() error {
	for _, file := range p.files {
		executor.Run("pip", "install", "-r", file)
	}
	return nil
}
