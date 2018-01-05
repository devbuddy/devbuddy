package features

import (
	"fmt"
	"path/filepath"

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

func (p Python) Enable(cfg *config.Config, proj *project.Project, env *Env, ui *termui.HookUI) error {
	name := fmt.Sprintf("%s-%s", proj.Slug(), p.Version)
	path := filepath.Join(cfg.DataDir, "virtualenvs", name)

	if !config.PathExists(path) {
		return DevUpNeeded
	}

	// - Add venv bin path to PATH

	env.Set("VIRTUAL_ENV", path)

	return nil
}

func (p Python) Disable(cfg *config.Config, proj *project.Project, env *Env, ui *termui.HookUI) {
	env.Unset("VIRTUAL_ENV")

	// - Remove virtualenv bin path from PATH
}
