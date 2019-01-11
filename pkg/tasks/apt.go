package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.RegisterTaskDefinition("apt", "Apt", parserApt).SetOsRequirement("debian")
}

func parserApt(config *taskapi.TaskConfig, task *taskapi.Task) error {
	packages, err := config.GetListOfStrings()
	if err != nil {
		return err
	}

	if len(packages) == 0 {
		return fmt.Errorf("no Apt packages specified")
	}

	task.Header = strings.Join(packages, ", ")

	task.AddAction(&aptInstall{packageNames: packages})

	return nil
}

type aptInstall struct {
	packageNames        []string
	missingPackageNames []string
}

func (a *aptInstall) Description() string {
	return ""
}

func (a *aptInstall) Needed(ctx *taskapi.Context) *taskapi.ActionResult {
	a.missingPackageNames = []string{}

	for _, name := range a.packageNames {
		result := shellSilent(ctx, fmt.Sprintf("dpkg -s \"%s\" | grep -q 'Status: install'", name)).Capture()
		if result.LaunchError != nil {
			return taskapi.ActionFailed("failed to check if package is installed: %s", result.LaunchError)
		}
		if result.Code != 0 {
			a.missingPackageNames = append(a.missingPackageNames, name)
		}
	}

	if len(a.missingPackageNames) > 0 {
		return taskapi.ActionNeeded("packages are not installed: %s", strings.Join(a.missingPackageNames, ", "))
	}

	return taskapi.ActionNotNeeded()
}

func (a *aptInstall) Run(ctx *taskapi.Context) error {
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
