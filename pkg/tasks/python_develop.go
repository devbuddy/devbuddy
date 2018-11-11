package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers/store"
)

func init() {
	t := registerTaskDefinition("python_develop")
	t.name = "Python develop"
	t.requiredTask = pythonTaskName
	t.parser = parserPythonDevelop
}

func parserPythonDevelop(config *taskConfig, task *Task) error {
	task.addAction(&pythonDevelopInstall{})
	return nil
}

type pythonDevelopInstall struct {
}

func (p *pythonDevelopInstall) description() string {
	return "install python package in develop mode"
}

func (p *pythonDevelopInstall) needed(ctx *context) (bool, error) {
	return store.New(ctx.proj.Path).HasFileChanged("setup.py")
}

func (p *pythonDevelopInstall) run(ctx *context) error {
	result := command(ctx, "pip", "install", "--require-virtualenv", "-e", ".").
		AddOutputFilter("already satisfied").Run()

	if result.Error != nil {
		return fmt.Errorf("Pip failed: %s", result.Error)
	}

	return store.New(ctx.proj.Path).RecordFileChange("setup.py")
}
