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
	var files []string

	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			files = append(files, v)
		} else {
			return fmt.Errorf("invalid pip files")
		}
	}
	if len(files) == 0 {
		return fmt.Errorf("no pip files specified")
	}

	task.header = strings.Join(files, ", ")

	for _, file := range files {
		action := newAction(
			fmt.Sprintf("install %s", file),
			func(ctx *context) error {
				result := command(ctx, "pip", "install", "--require-virtualenv", "-r", file).
					AddOutputFilter("already satisfied").Run()

				if result.Error != nil {
					return fmt.Errorf("Pip failed: %s", result.Error)
				}
				return nil
			})
		action.onFileChange(file).onFeatureChange("python")
		task.addAction(action)
	}

	return nil
}
