package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/executor"
)

type Custom struct {
	condition string
	command   string
}

func (c *Custom) Load(definition map[interface{}]interface{}) (bool, error) {
	if payload, ok := definition["custom"]; ok {
		properties := payload.(map[interface{}]interface{})

		command, _ := properties["meet"]
		condition, _ := properties["met?"]
		c.command = command.(string)
		c.condition = condition.(string)
		return true, nil
	}
	return false, nil
}

func (c *Custom) Perform() error {
	fmt.Printf("Task custom: command=\"%s\" condition=\"%s\"\n", c.command, c.condition)
	err := executor.RunShell(c.condition)
	if err != nil {
		err = executor.RunShell(c.command)
		if err != nil {
			fmt.Printf("Command failed: %s", err)
		}
	}
	return nil
}
