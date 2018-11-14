package tasks

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

func (c *customAction) needed(ctx *context) *actionResult {
	result := shellSilent(ctx, c.condition).Run()

	if result.LaunchError != nil {
		return actionFailed("failed to run the condition command: %s", result.LaunchError)
	}

	if result.Code != 0 {
		return actionNeeded("the met? command exited with a non-zero code")
	}
	return actionNotNeeded()
}

func (c *customAction) run(ctx *context) error {
	result := shell(ctx, c.command).Run()
	return result.Error
}
