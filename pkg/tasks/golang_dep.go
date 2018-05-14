package tasks

import (
	"fmt"
	"os/exec"
)

func init() {
	allTasks["golang_dep"] = newGolangDep
}

// GolangDep task manage the Go dependencies with Go Dep
type GolangDep struct {
}

func newGolangDep(config *taskConfig) (Task, error) {
	return &GolangDep{}, nil
}

func (d *GolangDep) name() string {
	return "Go Dep"
}

func (d *GolangDep) header() string {
	return "dep ensure"
}

func (d *GolangDep) preRunValidation(ctx *context) (err error) {
	_, hasFeature := ctx.features["golang"]
	if !hasFeature {
		return fmt.Errorf("You must specify a Go environment to use this task")
	}
	return nil
}

func (d *GolangDep) actions(ctx *context) []taskAction {
	return []taskAction{
		&golangDepInstall{},
		&golangDepEnsure{},
	}
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
	code, err := runCommand(ctx, "go", "get", "-u", "github.com/golang/dep/cmd/dep")
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("failed to install Go GolangDep. exit code: %d", code)
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
	code, err := runCommand(ctx, "dep", "ensure")
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("failed to run dep ensure. exit code: %d", code)
	}
	return nil
}
