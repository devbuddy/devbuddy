package tasks

import (
	"fmt"
	"os/exec"

	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("golang_dep", "Go Dep", parserGolangDep).SetRequiredTask("golang")
}

func parserGolangDep(config *taskapi.TaskConfig, task *taskapi.Task) error {
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

func (p *golangDepInstall) Needed(ctx *taskapi.Context) *taskapi.ActionResult {
	_, err := exec.LookPath("dep") // Just check if `dep` is in the PATH for now
	if err != nil {
		return taskapi.ActionNeeded("could not find the dep command in the PATH")
	}
	return taskapi.ActionNotNeeded()
}

func (p *golangDepInstall) Run(ctx *taskapi.Context) error {
	result := command(ctx, "go", "get", "-u", "github.com/golang/dep/cmd/dep").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to install Go GolangDep: %s", result.Error)
	}
	return nil
}

func (p *golangDepInstall) Feature() *features.FeatureInfo {
	return nil
}

type golangDepEnsure struct {
}

func (p *golangDepEnsure) Description() string {
	return "Run dep ensure"
}

func (p *golangDepEnsure) Needed(ctx *taskapi.Context) *taskapi.ActionResult {
	if !fileExists(ctx, "vendor") {
		return taskapi.ActionNeeded("the vendor directory does not exist")
	}

	// Is the vendor dir out dated?
	vendorMod, err := fileModTime(ctx, "vendor")
	if err != nil {
		return taskapi.ActionFailed("failed to get the modification of the vendor directory", err)
	}
	tomlMod, err := fileModTime(ctx, "Gopkg.toml")
	if err != nil {
		return taskapi.ActionFailed("failed to get the modification of Gopkg.toml", err)
	}
	lockMod, err := fileModTime(ctx, "Gopkg.lock")
	if err != nil {
		return taskapi.ActionFailed("failed to get the modification of Gopkg.lock", err)
	}
	if tomlMod > vendorMod || lockMod > vendorMod {
		return taskapi.ActionNeeded("Gopkg.toml or Gopkg.lock has been changed")
	}

	return taskapi.ActionNotNeeded()
}

func (p *golangDepEnsure) Run(ctx *taskapi.Context) error {
	result := command(ctx, "dep", "ensure").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to run dep ensure: %s", result.Error)
	}
	return nil
}

func (p *golangDepEnsure) Feature() *features.FeatureInfo {
	return nil
}
