package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
)

func init() {
	allTasks["python"] = NewPython
}

type Python struct {
	version string
}

func NewPython() Task {
	return &Python{}
}

func (p *Python) Load(definition interface{}) (bool, error) {
	def, ok := definition.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}
	if version, ok := def["python"]; ok {
		p.version, ok = version.(string)
		if !ok {
			return false, nil
		}
		return true, nil
	}

	return false, nil
}

func (p *Python) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader("Python", p.version)

	installed, err := p.InstallPython(ctx)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	venvInstalled, err := p.InstallVirtualEnv(ctx)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	venvCreated, err := p.CreateVirtualEnv(ctx)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	if installed || venvInstalled || venvCreated {
		ctx.ui.TaskActed()
	} else {
		ctx.ui.TaskAlreadyOk()
	}

	return nil
}

func (p *Python) InstallPython(ctx *Context) (acted bool, err error) {
	output, code, err := executor.Capture("pyenv", "versions", "--bare", "--skip-aliases")
	if err != nil {
		return
	}
	if code != 0 {
		return false, fmt.Errorf("failed to run pyenv versions. exit code: %d", code)
	}

	installedVersions := strings.Split(strings.TrimSpace(output), "\n")

	if stringInSlice(p.version, installedVersions) {
		return
	}

	code, err = executor.Run("pyenv", "install", p.version)
	if err != nil {
		return
	}
	if code != 0 {
		return false, fmt.Errorf("failed to install the required python version. exit code: %d", code)
	}

	return true, nil
}

func (p *Python) InstallVirtualEnv(ctx *Context) (acted bool, err error) {
	virtualenvExecutablePath := ctx.cfg.HomeDir(".pyenv", "versions", p.version, "bin", "virtualenv")
	pipExecutablePath := ctx.cfg.HomeDir(".pyenv", "versions", p.version, "bin", "pip")

	if config.PathExists(virtualenvExecutablePath) {
		return false, nil
	}

	code, err := executor.Run(pipExecutablePath, "install", "virtualenv")
	if err != nil {
		return
	}
	if code != 0 {
		return false, fmt.Errorf("failed to install virtualenv for the required python version. exit code: %d", code)
	}

	return true, nil
}

func (p *Python) CreateVirtualEnv(ctx *Context) (acted bool, err error) {
	virtualenvExecutablePath := ctx.cfg.HomeDir(".pyenv", "versions", p.version, "bin", "virtualenv")
	pythonExecutablePath := ctx.cfg.HomeDir(".pyenv", "versions", p.version, "bin", "python")

	name := fmt.Sprintf("%s-%s", ctx.proj.Slug(), p.version)
	virtualenvPath := ctx.cfg.DataDir("virtualenvs", name)

	if config.PathExists(virtualenvPath) {
		return false, nil
	}

	err = os.MkdirAll(filepath.Dir(virtualenvPath), 0750)
	if err != nil {
		return
	}

	code, err := executor.Run(virtualenvExecutablePath, "-p", pythonExecutablePath, virtualenvPath)
	if err != nil {
		return
	}
	if code != 0 {
		return false, fmt.Errorf("failed to create the virtualenv. exit code: %d", code)
	}

	return true, nil
}

func (p *Python) Features() map[string]string {
	return map[string]string{"python": p.version}
}
