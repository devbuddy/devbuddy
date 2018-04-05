package helpers

import (
	"fmt"
	"os"
	"path"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/utils"
)

type Golang struct {
	version string
	path    string
}

func NewGolang(cfg *config.Config, version string) *Golang {
	path := cfg.DataDir("golang", version)
	return &Golang{version: version, path: path}
}

func (g *Golang) Exists() bool {
	return utils.PathExists(g.path)
}

func (g *Golang) Path() string {
	return g.path
}

func (g *Golang) BinPath() string {
	return path.Join(g.path, "bin")
}

func (g *Golang) Install() (err error) {
	err = os.MkdirAll(g.path, 0750)
	if err != nil {
		return
	}

	os := "darwin"
	arch := "amd64"
	url := fmt.Sprintf("https://dl.google.com/go/go%s.%s-%s.tar.gz", g.version, os, arch)
	cmdline := fmt.Sprintf("curl -sL %s | tar --strip 1 -xzC %s", url, g.path)
	_, err = executor.RunShell(cmdline)

	return
}
