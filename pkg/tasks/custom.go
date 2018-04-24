package tasks

import (
	"fmt"
)

func init() {
	allTasks["custom"] = newCustom
}

type Custom struct {
	condition string
	command   string
}

func newCustom(config *taskConfig) (Task, error) {
	task := &Custom{}

	properties := config.payload.(map[interface{}]interface{})

	command, ok := properties["meet"]
	if !ok {
		return nil, fmt.Errorf("missing key 'meet'")
	}
	condition, ok := properties["met?"]
	if !ok {
		return nil, fmt.Errorf("missing key 'met?'")
	}

	var err error
	task.command, err = asString(command)
	if err != nil {
		return nil, fmt.Errorf("invalid meet value: %s", err)
	}
	task.condition, err = asString(condition)
	if err != nil {
		return nil, fmt.Errorf("invalid met? value: %s", err)
	}

	return task, nil
}

func (c *Custom) name() string {
	return "Custom"
}

func (c *Custom) header() string {
	return c.command
}

func (c *Custom) perform(ctx *Context) error {
	ran, err := c.runCommand(ctx)
	if err != nil {
		return err
	}

	if ran {
		ctx.ui.TaskActed()
	} else {
		ctx.ui.TaskAlreadyOk()
	}

	return nil
}

func (c *Custom) runCommand(ctx *Context) (bool, error) {
	code, err := runShellSilent(ctx, c.condition)
	if err != nil {
		return false, fmt.Errorf("failed to run the condition command: %s", err)
	}
	if code == 0 {
		return false, nil
	}

	// The condition command was run and returned a non-zero exit code.
	// It means we should run this custom task

	code, err = runShellSilent(ctx, c.command)
	if err != nil {
		return false, fmt.Errorf("command failed: %s", err)
	}
	if code != 0 {
		return false, fmt.Errorf("command exited with code %d", code)
	}

	return true, nil
}
