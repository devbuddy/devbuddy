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

	task.addAction(&aptInstall{packageNames: packages})

	return nil
}

type aptInstall struct {
	packageNames        []string
	missingPackageNames []string
}

func (a *aptInstall) description() string {
	return ""
}

func (a *aptInstall) needed(ctx *context) (bool, error) {
	a.missingPackageNames = []string{}

	for _, name := range a.packageNames {
		result := shellSilent(ctx, fmt.Sprintf("dpkg -s \"%s\" | grep -q 'Status: install'", name)).Capture()
		if result.LaunchError != nil {
			return false, fmt.Errorf("failed to check if package is installed: %s", result.LaunchError)
		}
		if result.Code != 0 {
			a.missingPackageNames = append(a.missingPackageNames, name)
		}
	}

	if len(a.missingPackageNames) > 0 {
		return true, fmt.Errorf("packages are not installed: %s", strings.Join(a.missingPackageNames, ", "))
	}
	return false, nil
}

func (a *aptInstall) run(ctx *context) error {
	result := sudoCommand(ctx, "apt-get", "update").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to run apt-get update: %s", result.Error)
	}

	args := append([]string{"install", "--no-install-recommends", "-y"}, a.missingPackageNames...)
	result = sudoCommand(ctx, "apt-get", args...).Run()
	if result.Error != nil {
		return fmt.Errorf("failed to run apt-get install: %s", result.Error)
	}

	return nil
}
