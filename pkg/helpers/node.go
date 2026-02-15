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

type Node struct {
	ctx     *context.Context
	version string
	path    string
	tarDir  string
}

func NewNode(ctx *context.Context, version string) *Node {
	return &Node{
		ctx:     ctx,
		version: version,
		path:    ctx.Cfg.DataDir("node", version),
		tarDir:  ctx.Cfg.DataDir("node"),
	}
}

func (n *Node) Exists() bool {
	return utils.PathExists(n.Which("node"))
}

func (n *Node) Path() string {
	return n.path
}

func (n *Node) BinPath() string {
	return path.Join(n.path, "bin")
}

func (n *Node) Which(program string) string {
	return path.Join(n.path, "bin", program)
}

func NodeJSIsSupportedOnThisPlatform() bool {
	return runtime.GOARCH == "amd64"
}

func (n *Node) Install() error {
	if !NodeJSIsSupportedOnThisPlatform() {
		return fmt.Errorf("NodeJS installation is only supported on %s by DevBuddy", runtime.GOARCH)
	}

	archiveName := fmt.Sprintf("node-v%s-%s-x64.tar.gz", n.version, runtime.GOOS)
	tarPath := path.Join(n.tarDir, archiveName)

	if !utils.PathExists(tarPath) {
		err := os.MkdirAll(n.tarDir, 0750)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("https://nodejs.org/dist/v%s/%s", n.version, archiveName)
		err = NewDownloader(url).DownloadToFile(tarPath)
		if err != nil {
			return fmt.Errorf("failed to download NodeJS %s from %s: %w", n.version, url, err)
		}
	}

	err := os.MkdirAll(n.path, 0750)
	if err != nil {
		return err
	}

	result := n.ctx.Executor.Run(executor.New("tar", "--strip", "1", "-xzC", n.path, "-f", tarPath))
	if result.Error != nil {
		return fmt.Errorf("failed to extract %s to %s: %w", tarPath, n.path, result.Error)
	}

	return nil
}
