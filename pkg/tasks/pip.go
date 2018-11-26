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
		// task.addAction(&pipInstall{file: file})
		task.
			addActionWithBuilder(fmt.Sprintf("install %s", file),
				func(ctx *context) error {
					result := command(ctx, "pip", "install", "--require-virtualenv", "-r", file).
						AddOutputFilter("already satisfied").Run()

					if result.Error != nil {
						return fmt.Errorf("Pip failed: %s", result.Error)
					}
					return nil
				}).
			addFileChangeCondition(file).
			addFeatureChangeCondition("python")
	}

	return nil
}

type pipInstall struct {
	file    string
	success bool
}

func (p *pipInstall) description() string {
	return fmt.Sprintf("install %s", p.file)
}

func (p *pipInstall) needed(ctx *context) *actionResult {
	if !p.success {
		return actionNeeded("")
	}
	return actionNotNeeded()
}

func (p *pipInstall) run(ctx *context) error {
	result := command(ctx, "pip", "install", "--require-virtualenv", "-r", p.file).
		AddOutputFilter("already satisfied").Run()

	if result.Error != nil {
		return fmt.Errorf("Pip failed: %s", result.Error)
	}
	p.success = true
	return nil
}
