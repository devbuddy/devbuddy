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

func (p *golangDepInstall) needed(ctx *context) (bool, error) {
	_, err := exec.LookPath("dep") // Just check if `dep` is in the PATH for now
	return err != nil, nil
}

func (p *golangDepInstall) run(ctx *context) error {
	err := command(ctx, "go", "get", "-u", "github.com/golang/dep/cmd/dep").Run()
	if err != nil {
		return fmt.Errorf("failed to install Go GolangDep: %s", err)
	}
	return nil
}

type golangDepEnsure struct {
}

func (p *golangDepEnsure) description() string {
	return "Run dep ensure"
}

func (p *golangDepEnsure) needed(ctx *context) (bool, error) {
	if !fileExists(ctx, "vendor") {
		return true, nil
	}

	// Is the vendor dir out dated?
	vendorMod, err := fileModTime(ctx, "vendor")
	if err != nil {
		return false, err
	}
	tomlMod, err := fileModTime(ctx, "Gopkg.toml")
	if err != nil {
		return false, err
	}
	lockMod, err := fileModTime(ctx, "Gopkg.lock")
	if err != nil {
		return false, err
	}
	if tomlMod > vendorMod || lockMod > vendorMod {
		return true, nil
	}

	return false, nil
}

func (p *golangDepEnsure) run(ctx *context) error {
	err := command(ctx, "dep", "ensure").Run()
	if err != nil {
		return fmt.Errorf("failed to run dep ensure: %s", err)
	}
	return nil
}
