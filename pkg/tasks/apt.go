package tasks

import (
	"fmt"
	"strings"
)

func init() {
	t := registerTaskDefinition("apt")
	t.name = "Apt"
	t.parser = parserApt
}

func parserApt(config *taskConfig, task *Task) error {
	packages, err := config.getListOfStrings()
	if err != nil {
		return err
	}

	if len(packages) == 0 {
		return fmt.Errorf("no Apt packages specified")
	}

	task.header = strings.Join(packages, ", ")

	for _, p := range packages {
		task.addAction(&aptInstall{packageName: p})
	}

	return nil
}

type aptInstall struct {
	packageName string
}

func (a *aptInstall) description() string {
	return fmt.Sprintf("installing %s", a.packageName)
}

func (a *aptInstall) needed(ctx *context) (bool, error) {
	result := commandSilent(ctx, "dpkg", "-s", a.packageName).Capture()
	if result.Error != nil {
		return false, fmt.Errorf("failed to check if package is installed: %s", result.Error)
	}
	return result.Code != 0, nil
}

func (a *aptInstall) run(ctx *context) error {
	result := command(ctx, "apt", "install", a.packageName).Run()

	if result.Error != nil {
		return fmt.Errorf("Apt failed: %s", result.Error)
	}

	return nil
}
