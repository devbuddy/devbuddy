package tasks

import (
	"fmt"
	"strings"
)

func init() {
	Register("apt", "Apt", parserApt).SetOsRequirement("debian")
}

func parserApt(config *TaskConfig, task *Task) error {
	packages, err := config.getListOfStrings()
	if err != nil {
		return err
	}

	if len(packages) == 0 {
		return fmt.Errorf("no Apt packages specified")
	}

	task.SetInfo(strings.Join(packages, ", "))

	task.AddAction(&aptInstall{packageNames: packages})

	return nil
}

type aptInstall struct {
	packageNames        []string
	missingPackageNames []string
}

func (a *aptInstall) description() string {
	return ""
}

func (a *aptInstall) needed(ctx *Context) *actionResult {
	a.missingPackageNames = []string{}

	for _, name := range a.packageNames {
		result := shellSilent(ctx, fmt.Sprintf("dpkg -s \"%s\" | grep -q 'Status: install'", name)).Capture()
		if result.LaunchError != nil {
			return actionFailed("failed to check if package is installed: %s", result.LaunchError)
		}
		if result.Code != 0 {
			a.missingPackageNames = append(a.missingPackageNames, name)
		}
	}

	if len(a.missingPackageNames) > 0 {
		return actionNeeded("packages are not installed: %s", strings.Join(a.missingPackageNames, ", "))
	}

	return actionNotNeeded()
}

func (a *aptInstall) run(ctx *Context) error {
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
