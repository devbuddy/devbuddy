package tasks

import (
	"fmt"
	"strings"
)

func init() {
	Register("python_develop", "Python develop", parserPythonDevelop).SetRequiredTask(pythonTaskName)
}

func parserPythonDevelop(config *TaskConfig, task *Task) error {
	extras, err := config.GetListOfStringsPropertyDefault("extras", []string{})
	if err != nil {
		return err
	}

	pipTarget := "."
	if len(extras) > 0 {
		pipTarget = fmt.Sprintf(".[%s]", strings.Join(extras, ","))
	}
	pipArgs := []string{"install", "--require-virtualenv", "-e", pipTarget}

	builder := actionBuilder("install python package in develop mode", func(ctx *Context) error {
		result := command(ctx, "pip", pipArgs...).AddOutputFilter("already satisfied").Run()
		if result.Error != nil {
			return fmt.Errorf("Pip failed: %s", result.Error)
		}

		return nil
	})
	builder.OnFileChange("setup.py")
	builder.OnFileChange("setup.cfg")

	task.AddAction(builder.Build())
	return nil
}
