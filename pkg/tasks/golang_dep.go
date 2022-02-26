package tasks

import (
	"fmt"
	"os/exec"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("golang_dep", "Go Dep", parserGolangDep).SetRequiredTask("go")
}

func parserGolangDep(config *api.TaskConfig, task *api.Task) error {
	task.Info = "dep ensure"
	task.AddAction(&golangDepInstall{})
	task.AddAction(&golangDepEnsure{})
	return nil
}

type golangDepInstall struct {
}

func (p *golangDepInstall) Description() string {
	return "Install Go Dep"
}

func (p *golangDepInstall) Needed(ctx *context.Context) *api.ActionResult {
	_, err := exec.LookPath("dep") // Just check if `dep` is in the PATH for now
	if err != nil {
		return api.Needed("could not find the dep command in the PATH")
	}
	return api.NotNeeded()
}

func (p *golangDepInstall) Run(ctx *context.Context) error {
	result := command(ctx, "go", "get", "-u", "github.com/golang/dep/cmd/dep").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to install Go GolangDep: %w", result.Error)
	}
	return nil
}

func (p *golangDepInstall) Feature() *autoenv.FeatureInfo {
	return nil
}

type golangDepEnsure struct {
}

func (p *golangDepEnsure) Description() string {
	return "Run dep ensure"
}

func (p *golangDepEnsure) Needed(ctx *context.Context) *api.ActionResult {
	if !fileExists(ctx, "vendor") {
		return api.Needed("the vendor directory does not exist")
	}

	// Is the vendor dir out dated?
	vendorMod, err := fileModTime(ctx, "vendor")
	if err != nil {
		return api.Failed("failed to get the modification of the vendor directory: %s", err)
	}
	tomlMod, err := fileModTime(ctx, "Gopkg.toml")
	if err != nil {
		return api.Failed("failed to get the modification of Gopkg.toml: %s", err)
	}
	lockMod, err := fileModTime(ctx, "Gopkg.lock")
	if err != nil {
		return api.Failed("failed to get the modification of Gopkg.lock: %s", err)
	}
	if tomlMod > vendorMod || lockMod > vendorMod {
		return api.Needed("Gopkg.toml or Gopkg.lock has been changed")
	}

	return api.NotNeeded()
}

func (p *golangDepEnsure) Run(ctx *context.Context) error {
	result := command(ctx, "dep", "ensure").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to run dep ensure: %w", result.Error)
	}
	return nil
}

func (p *golangDepEnsure) Feature() *autoenv.FeatureInfo {
	return nil
}
