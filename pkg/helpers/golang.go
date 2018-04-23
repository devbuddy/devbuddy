package helpers

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/utils"
)

type Golang struct {
	version string
	path    string
	tarDir  string
}

func NewGolang(cfg *config.Config, version string) *Golang {
	return &Golang{
		version: version,
		path:    cfg.DataDir("golang", version),
		tarDir:  cfg.DataDir("golang"),
	}
}

func (g *Golang) Exists() bool {
	return utils.PathExists(g.Which("go"))
}

func (g *Golang) Path() string {
	return g.path
}

func (g *Golang) BinPath() string {
	return path.Join(g.path, "bin")
}

func (g *Golang) Which(program string) string {
	return path.Join(g.path, "bin", program)
}

func (g *Golang) Install() (err error) {
	archiveName := fmt.Sprintf("go%s.%s-%s.tar.gz", g.version, runtime.GOOS, runtime.GOARCH)
	tarPath := path.Join(g.tarDir, archiveName)

	if !utils.PathExists(tarPath) {
		err = os.MkdirAll(g.tarDir, 0750)
		if err != nil {
			return
		}

		url := "https://dl.google.com/go/" + archiveName
		err = utils.DownloadFile(tarPath, url)
		if err != nil {
			return fmt.Errorf("failed to download Go %s from %s: %s", g.version, url, err)
		}
	}

	err = os.MkdirAll(g.path, 0750)
	if err != nil {
		return
	}

	code, err := executor.New("tar", "--strip", "1", "-xzC", g.path, "-f", tarPath).Run()
	if err != nil || code != 0 {
		return fmt.Errorf("failed to extract %s to %s (code: %d err: %s)", tarPath, g.path, code, err)
	}

	return nil
}
