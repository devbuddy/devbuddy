package tasks

import (
	"fmt"
	"strings"
)

func init() {
	t := registerTaskDefinition("pip")
	t.name = "Pip"
	t.requiredTask = pythonTaskName
	t.parser = parserPip
}

func parserPip(config *taskConfig, task *Task) error {
	files, err := config.getListOfStrings()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no pip files specified")
	}

	task.header = strings.Join(files, ", ")

	for _, file := range files {
		builder := actionBuilder(fmt.Sprintf("install %s", file), func(ctx *Context) error {
			pipArgs := []string{"install", "--require-virtualenv", "-r", file}
			result := command(ctx, "pip", pipArgs...).AddOutputFilter("already satisfied").Run()
			if result.Error != nil {
				return fmt.Errorf("Pip failed: %s", result.Error)
			}
			return nil
		})
		builder.OnFileChange(file)
		task.addAction(builder.Build())
	}

	return nil
}
