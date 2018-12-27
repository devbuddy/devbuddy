package tasks

func init() {
	t := registerTaskDefinition("custom")
	t.name = "Custom"
	t.parser = parserCustom
}

func parserCustom(config *taskConfig, task *Task) error {
	command, err := config.getStringProperty("meet")
	if err != nil {
		return err
	}
	condition, err := config.getStringProperty("met?")
	if err != nil {
		return err
	}
	name, err := config.getStringPropertyDefault("name", command)
	if err != nil {
		return err
	}

	task.header = name

	builder := actionBuilder("", func(ctx *Context) error {
		result := shell(ctx, command).Run()
		return result.Error
	})
	builder.OnFunc(func(ctx *Context) *actionResult {
		result := shellSilent(ctx, condition).Capture()
		if result.LaunchError != nil {
			return actionFailed("failed to run the condition command: %s", result.LaunchError)
		}
		if result.Code != 0 {
			return actionNeeded("the met? command exited with a non-zero code")
		}
		return actionNotNeeded()
	})
	task.addAction(builder.Build())

	return nil
}
