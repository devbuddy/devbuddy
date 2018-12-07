package tasks

import (
	"fmt"
)

func init() {
	t := registerTaskDefinition("python_develop")
	t.name = "Python develop"
	t.requiredTask = pythonTaskName
	t.parser = parserPythonDevelop
}

func parserPythonDevelop(config *taskConfig, task *Task) error {
	builder := actionBuilder("install python package in develop mode", func(ctx *context) error {
		result := command(ctx, "pip", "install", "--require-virtualenv", "-e", ".").
			AddOutputFilter("already satisfied").Run()

		if result.Error != nil {
			return fmt.Errorf("Pip failed: %s", result.Error)
		}

		return nil
	})
	builder.OnFileChange("setup.py")

	task.addAction(builder.Build())
	return nil
}
