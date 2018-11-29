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

	action := newAction("", func(ctx *context) error {
		missingPackages, err := findMissingAptPackages(ctx, packages)
		if err != nil {
			return err
		}

		result := sudoCommand(ctx, "apt-get", "update").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to run apt-get update: %s", result.Error)
		}

		args := append([]string{"install", "--no-install-recommends", "-y"}, missingPackages...)
		result = sudoCommand(ctx, "apt-get", args...).Run()
		if result.Error != nil {
			return fmt.Errorf("failed to run apt-get install: %s", result.Error)
		}

		return nil
	}).onFunc(func(ctx *context) *actionResult {
		missingPackages, err := findMissingAptPackages(ctx, packages)
		if err != nil {
			return actionFailed("failed to check if package is installed: %s", err)
		}

		if len(missingPackages) > 0 {
			return actionNeeded("packages are not installed: %s", strings.Join(missingPackages, ", "))
		}

		return actionNotNeeded()
	})
	task.addAction(action)

	return nil
}

func findMissingAptPackages(ctx *context, packages []string) ([]string, error) {
	missing := []string{}
	for _, name := range packages {
		result := shellSilent(ctx, fmt.Sprintf("dpkg -s \"%s\" | grep -q 'Status: install'", name)).Capture()
		if result.LaunchError != nil {
			return nil, result.LaunchError
		}
		if result.Code != 0 {
			missing = append(missing, name)
		}
	}
	return missing, nil
}
