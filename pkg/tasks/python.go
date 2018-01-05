package tasks

import (
	"fmt"
	"strings"

	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/termui"
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

func (p *Python) Perform(ui *termui.UI) (err error) {
	ui.TaskHeader("Python", p.version)

	installed, err := p.InstallPython(ui)
	if err != nil {
		ui.TaskError(err)
		return err
	}
	created, err := p.CreateVirtualEnv(ui)
	if err != nil {
		ui.TaskError(err)
		return err
	}
	if installed || created {
		ui.TaskActed()
	} else {
		ui.TaskAlreadyOk()
	}

	return nil
}

func (p *Python) InstallPython(ui *termui.UI) (acted bool, err error) {
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

func (p *Python) CreateVirtualEnv(ui *termui.UI) (acted bool, err error) {
	return false, nil
}

func (p *Python) Features() map[string]string {
	return map[string]string{"python": p.version}
}
