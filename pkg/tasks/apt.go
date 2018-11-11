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
	code, err := commandSilent(ctx, "dpkg", "-s", a.packageName).RunWithCode()
	if err != nil {
		return false, fmt.Errorf("failed to check if package is installed: %s", err)
	}
	return code != 0, nil
}

func (a *aptInstall) run(ctx *context) error {
	err := command(ctx, "apt", "install", a.packageName).Run()

	if err != nil {
		return fmt.Errorf("Apt failed: %s", err)
	}

	return nil
}
