package tasks

func init() {
	Register("custom", "Custom", parserCustom)
}

func parserCustom(config *TaskConfig, task *Task) error {
	command, err := config.GetStringProperty("meet")
	if err != nil {
		return err
	}
	condition, err := config.GetStringProperty("met?")
	if err != nil {
		return err
	}
	name, err := config.GetStringPropertyDefault("name", command)
	if err != nil {
		return err
	}

	task.SetInfo(name)

	builder := actionBuilder("", func(ctx *Context) error {
		result := shell(ctx, command).Run()
		return result.Error
	})
	builder.OnFunc(func(ctx *Context) *ActionResult {
		result := shellSilent(ctx, condition).Capture()
		if result.LaunchError != nil {
			return ActionFailed("failed to run the condition command: %s", result.LaunchError)
		}
		if result.Code != 0 {
			return ActionNeeded("the met? command exited with a non-zero code")
		}
		return ActionNotNeeded()
	})
	task.AddAction(builder.Build())

	return nil
}
