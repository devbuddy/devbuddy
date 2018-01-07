package helpers

import (
	"fmt"
	"path"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
)

type Virtualenv struct {
	name string
	path string
}

func NewVirtualenv(cfg *config.Config, proj *project.Project, pythonVersion string) *Virtualenv {
	name := fmt.Sprintf("%s-%s", proj.Slug(), pythonVersion)
	path := cfg.DataDir("virtualenvs", name)

	v := Virtualenv{name: name, path: path}
	return &v
}

func (v *Virtualenv) Exists() bool {
	return config.PathExists(v.path)
}

func (v *Virtualenv) Name() string {
	return v.name
}

func (v *Virtualenv) Path() string {
	return v.path
}

func (v *Virtualenv) BinPath() string {
	return path.Join(v.path, "bin")
}
