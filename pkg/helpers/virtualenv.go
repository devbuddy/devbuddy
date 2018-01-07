package helpers

import (
	"fmt"
	"path"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
)

type Virtualenv struct {
	path string
}

func NewVirtualenv(cfg *config.Config, name string) *Virtualenv {
	path := cfg.DataDir("virtualenvs", name)
	v := Virtualenv{path: path}
	return &v
}

func (v *Virtualenv) Exists() bool {
	return config.PathExists(v.path)
}

func (v *Virtualenv) Path() string {
	return v.path
}

func (v *Virtualenv) BinPath() string {
	return path.Join(v.path, "bin")
}

func VirtualenvName(proj *project.Project, version string) string {
	return fmt.Sprintf("%s-%s", proj.Slug(), version)
}
