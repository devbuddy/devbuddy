package tasks

func init() {
	Register("custom", "Custom", parserCustom)
}

func parserCustom(config *TaskConfig, task *Task) error {
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

	task.SetInfo(name)

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
	task.AddAction(builder.Build())

	return nil
}
