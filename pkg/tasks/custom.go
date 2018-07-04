package tasks

import (
	"fmt"
)

func init() {
	t := registerTask("custom")
	t.name = "Custom"
	t.parser = parserCustom
}

func parserCustom(config *taskConfig, task *Task) error {
	properties := config.payload.(map[interface{}]interface{})

	name, ok := properties["name"]
	if !ok {
		name = ""
	}
	command, ok := properties["meet"]
	if !ok {
		return fmt.Errorf("missing key 'meet'")
	}
	condition, ok := properties["met?"]
	if !ok {
		return fmt.Errorf("missing key 'met?'")
	}

	nameStr, err := asString(name)
	if err != nil {
		return fmt.Errorf("invalid name value: %s", err)
	}
	commandStr, err := asString(command)
	if err != nil {
		return fmt.Errorf("invalid meet value: %s", err)
	}
	conditionStr, err := asString(condition)
	if err != nil {
		return fmt.Errorf("invalid met? value: %s", err)
	}

	if nameStr == "" {
		nameStr = commandStr
	}

	task.header = nameStr
	task.addAction(&customAction{condition: conditionStr, command: commandStr})

	return nil
}

type customAction struct {
	condition string
	command   string
}

func (c *customAction) description() string {
	return ""
}

func (c *customAction) needed(ctx *context) (bool, error) {
	code, err := shellSilent(ctx, c.condition).RunWithCode()
	if err != nil {
		return false, fmt.Errorf("failed to run the condition command: %s", err)
	}
	return code != 0, nil
}

func (c *customAction) run(ctx *context) error {
	err := shell(ctx, c.command).Run()
	if err != nil {
		return fmt.Errorf("command failed: %s", err)
	}
	return nil
}
