package features

import (
	"fmt"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

func init() {
	allFeatures["python"] = NewPython
}

type Python struct {
	Version string
}

func NewPython(param string) Feature {
	return Python{Version: param}
}

func (p Python) Enable(proj *project.Project, env *Env, ui *termui.UI) error {
	path := fmt.Sprintf("~/.pyenv/virtualenvs/%s-%s", proj.Slug(), p.Version)
	path = config.ExpandDir(path)

	if !config.PathExists(path) {
		return DevUpNeeded
	}

	// - Add venv bin path to PATH

	env.Set("VIRTUAL_ENV", path)

	return nil
}

func (p Python) Disable(proj *project.Project, env *Env, ui *termui.UI) {
	env.Unset("VIRTUAL_ENV")

	// - Remove virtualenv bin path from PATH
}
