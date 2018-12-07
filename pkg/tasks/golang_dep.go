package tasks

import (
	"fmt"
	"os/exec"
)

func init() {
	t := registerTaskDefinition("golang_dep")
	t.name = "Go Dep"
	t.requiredTask = "golang"
	t.parser = parserGolangDep
}

func parserGolangDep(config *taskConfig, task *Task) error {
	task.header = "dep ensure"
	task.addAction(&golangDepInstall{})
	task.addAction(&golangDepEnsure{})
	return nil
}

type golangDepInstall struct {
}

func (p *golangDepInstall) description() string {
	return "Install Go Dep"
}

func (p *golangDepInstall) needed(ctx *Context) *actionResult {
	_, err := exec.LookPath("dep") // Just check if `dep` is in the PATH for now
	if err != nil {
		return actionNeeded("could not find the dep command in the PATH")
	}
	return actionNotNeeded()
}

func (p *golangDepInstall) run(ctx *Context) error {
	result := command(ctx, "go", "get", "-u", "github.com/golang/dep/cmd/dep").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to install Go GolangDep: %s", result.Error)
	}
	return nil
}

type golangDepEnsure struct {
}

func (p *golangDepEnsure) description() string {
	return "Run dep ensure"
}

func (p *golangDepEnsure) needed(ctx *Context) *actionResult {
	if !fileExists(ctx, "vendor") {
		return actionNeeded("the vendor directory does not exist")
	}

	// Is the vendor dir out dated?
	vendorMod, err := fileModTime(ctx, "vendor")
	if err != nil {
		return actionFailed("failed to get the modification of the vendor directory", err)
	}
	tomlMod, err := fileModTime(ctx, "Gopkg.toml")
	if err != nil {
		return actionFailed("failed to get the modification of Gopkg.toml", err)
	}
	lockMod, err := fileModTime(ctx, "Gopkg.lock")
	if err != nil {
		return actionFailed("failed to get the modification of Gopkg.lock", err)
	}
	if tomlMod > vendorMod || lockMod > vendorMod {
		return actionNeeded("Gopkg.toml or Gopkg.lock has been changed")
	}

	return actionNotNeeded()
}

func (p *golangDepEnsure) run(ctx *Context) error {
	result := command(ctx, "dep", "ensure").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to run dep ensure: %s", result.Error)
	}
	return nil
}
