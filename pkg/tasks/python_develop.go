package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers/store"
)

func init() {
	t := registerTask("python_develop")
	t.name = "Python develop"
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
	err := command(ctx, "pip", "install", "--require-virtualenv", "-e", ".").AddOutputFilter("already satisfied").Run()
	if err != nil {
		return fmt.Errorf("Pip failed: %s", err)
	}

	return store.New(ctx.proj.Path).RecordFileChange("setup.py")
}
