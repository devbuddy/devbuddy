package tasks

import (
	"fmt"
)

func init() {
	t := registerTaskDefinition("custom")
	t.name = "Custom"
	t.parser = parserCustom
}

func parserCustom(config *taskConfig, task *Task) error {
	command, err := config.getStringProperty("meet", false)
	if err != nil {
		return err
	}
	condition, err := config.getStringProperty("met?", false)
	if err != nil {
		return err
	}
	name, err := config.getStringPropertyDefault("name", command, false)
	if err != nil {
		return err
	}

	task.header = name
	task.addAction(&customAction{condition: condition, command: command})

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
