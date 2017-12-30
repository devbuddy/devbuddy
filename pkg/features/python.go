package features

import (
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

func init() {
	FeatureMap["python"] = NewPython
}

type Python struct {
	Version string
}

func NewPython(param string) Feature {
	return Python{Version: param}
}

func (p Python) Enable(proj *project.Project, env *Env, ui *termui.UI) {
	// ui.Debug("Initial VIRTUAL_ENV=%s", env.Get("VIRTUAL_ENV"))

	// - Elect a virtualenv path (following dev choice of ~/.pyenv/virtualenvs/... ?)
	// - Warn and exit if virtualenv doesn't exist
	// - Add venv bin path to PATH
	// - Set the VIRTUAL_ENV variable
}

func (p Python) Disable(proj *project.Project, env *Env, ui *termui.UI) {
	// - Unset VIRTUAL_ENV
	// - Remove virtualenv bin path from PATH
}
