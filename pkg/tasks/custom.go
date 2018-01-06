package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/executor"
)

func init() {
	allTasks["custom"] = NewCustom
}

type Custom struct {
	condition string
	command   string
}

func NewCustom() Task {
	return &Custom{}
}

func (c *Custom) Load(definition interface{}) (bool, error) {
	def, ok := definition.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}

	if payload, ok := def["custom"]; ok {
		properties := payload.(map[interface{}]interface{})

		command, ok := properties["meet"]
		if !ok {
			return false, nil
		}
		condition, ok := properties["met?"]
		if !ok {
			return false, nil
		}
		c.command = command.(string)
		c.condition = condition.(string)
		return true, nil
	}
	return false, nil
}

func (c *Custom) Perform(ctx *Context) error {
	ctx.ui.TaskHeader("Custom", c.command)

	ran, err := c.runCommand()
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	if ran {
		ctx.ui.TaskActed()
	} else {
		ctx.ui.TaskAlreadyOk()
	}

	return nil
}

func (c *Custom) runCommand() (bool, error) {
	code, err := executor.RunShellSilent(c.condition)
	if err != nil {
		return false, fmt.Errorf("failed to run the condition command: %s", err)
	}
	if code == 0 {
		return false, nil
	}

	// The condition command was run and returned a non-zero exit code.
	// It means we should run this custom task

	code, err = executor.RunShellSilent(c.command)
	if err != nil {
		return false, fmt.Errorf("command failed: %s", err)
	}
	if code != 0 {
		return false, fmt.Errorf("command exited with code %d", code)
	}

	return true, nil
}
