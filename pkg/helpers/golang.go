package helpers

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/utils"
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

	err = executor.New("tar", "--strip", "1", "-xzC", g.path, "-f", tarPath).Run()
	if err != nil {
		return fmt.Errorf("failed to extract %s to %s: %s", tarPath, g.path, err)
	}

	return nil
}
