package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("homebrew", "Homebrew", parserHomebrew).SetOSRequirement("macos")
}

func parserHomebrew(config *api.TaskConfig, task *api.Task) error {
	formulas, err := config.GetListOfStrings()
	if err != nil {
		return err
	}

	if len(formulas) == 0 {
		return fmt.Errorf("no homebrew formulas specified")
	}

	task.Info = strings.Join(formulas, ", ")

	for _, f := range formulas {
		task.AddAction(&brewInstall{formula: f})
	}

	return nil
}

type brewInstall struct {
	formula string
}

func (b *brewInstall) Description() string {
	return fmt.Sprintf("installing %s", b.formula)
}

func (b *brewInstall) Needed(ctx *context.Context) *api.ActionResult {
	if helpers.HomeBrewPackageIsInstalled(b.formula) {
		return api.NotNeeded()
	}
	return api.Needed("package %s is not installed", b.formula)
}

func (b *brewInstall) Run(ctx *context.Context) error {
	ctx.UI.TaskCommand("brew", "install", b.formula)
	result := ctx.Executor.Run(executor.New("brew", "install", b.formula).AddEnvVar("HOMEBREW_NO_AUTO_UPDATE", "1"))
	if result.Error != nil {
		return fmt.Errorf("failed to run brew install: %w", result.Error)
	}

	return nil
}

func (b *brewInstall) Feature() *autoenv.FeatureInfo {
	return nil
}
