package helpers

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type Golang struct {
	ctx     *context.Context
	version string
	path    string
	tarDir  string
}

func NewGolang(ctx *context.Context, version string) *Golang {
	return &Golang{
		ctx:     ctx,
		version: version,
		path:    ctx.Cfg.DataDir("golang", version),
		tarDir:  ctx.Cfg.DataDir("golang"),
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
		err = NewDownloader(url).DownloadToFile(tarPath)
		if err != nil {
			return fmt.Errorf("failed to download Go %s from %s: %w", g.version, url, err)
		}
	}

	err = os.MkdirAll(g.path, 0750)
	if err != nil {
		return
	}

	result := g.ctx.Executor.Run(executor.New("tar", "--strip", "1", "-xzC", g.path, "-f", tarPath))
	if result.Error != nil {
		return fmt.Errorf("failed to extract %s to %s: %w", tarPath, g.path, result.Error)
	}

	return nil
}
