package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers/store"
)

func init() {
	allTasks["python_develop"] = newPythonDevelop
}

type pythonDevelop struct {
}

func newPythonDevelop(config *taskConfig) (Task, error) {
	return &pythonDevelop{}, nil
}

func (p *pythonDevelop) name() string {
	return "Python develop"
}

func (p *pythonDevelop) header() string {
	return ""
}

func (p *pythonDevelop) preRunValidation(ctx *context) (err error) {
	_, hasPythonFeature := ctx.features["python"]
	if !hasPythonFeature {
		return fmt.Errorf("You must specify a Python environment to use this task")
	}
	return nil
}

func (p *pythonDevelop) actions(ctx *context) (actions []taskAction) {
	actions = append(actions, &pythonDevelopInstall{})
	return
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
