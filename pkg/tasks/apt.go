package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("apt", "Apt", parserApt).SetOSRequirement("debian")
}

func parserApt(config *api.TaskConfig, task *api.Task) error {
	packages, err := config.GetListOfStrings()
	if err != nil {
		return err
	}

	if len(packages) == 0 {
		return fmt.Errorf("no Apt packages specified")
	}

	task.Info = strings.Join(packages, ", ")

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

func (a *aptInstall) Needed(ctx *context.Context) *api.ActionResult {
	a.missingPackageNames = []string{}

	for _, name := range a.packageNames {
		result := shellSilent(ctx, fmt.Sprintf("dpkg -s \"%s\" | grep -q 'Status: install'", name)).Capture()
		if result.LaunchError != nil {
			return api.Failed("failed to check if package is installed: %s", result.LaunchError)
		}
		if result.Code != 0 {
			a.missingPackageNames = append(a.missingPackageNames, name)
		}
	}

	if len(a.missingPackageNames) > 0 {
		return api.Needed("packages are not installed: %s", strings.Join(a.missingPackageNames, ", "))
	}

	return api.NotNeeded()
}

func (a *aptInstall) Run(ctx *context.Context) error {
	result := sudoCommand(ctx, "apt-get", "update").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to run apt-get update: %w", result.Error)
	}

	args := append([]string{"install", "--no-install-recommends", "-y"}, a.missingPackageNames...)
	result = sudoCommand(ctx, "apt-get", args...).Run()
	if result.Error != nil {
		return fmt.Errorf("failed to run apt-get install: %w", result.Error)
	}

	return nil
}

func (a *aptInstall) Feature() *autoenv.FeatureInfo {
	return nil
}
