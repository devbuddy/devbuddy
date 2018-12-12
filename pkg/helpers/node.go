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

type Node struct {
	version string
	path    string
	tarDir  string
}

func NewNode(cfg *config.Config, version string) *Node {
	return &Node{
		version: version,
		path:    cfg.DataDir("node", version),
		tarDir:  cfg.DataDir("node"),
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

func (n *Node) Install() error {
	if runtime.GOARCH != "amd64" {
		return fmt.Errorf("Binary distribution are not available for %s architecture", runtime.GOARCH)
	}

	archiveName := fmt.Sprintf("node-%s-%s-x64.tar.gz", n.version, runtime.GOOS)
	tarPath := path.Join(n.tarDir, archiveName)

	if !utils.PathExists(tarPath) {
		err := os.MkdirAll(n.tarDir, 0750)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("https://nodejs.org/dist/%s/%s", n.version, archiveName)
		err = utils.DownloadFile(tarPath, url)
		if err != nil {
			return fmt.Errorf("failed to download NodeJS %s from %s: %s", n.version, url, err)
		}
	}

	err := os.MkdirAll(n.path, 0750)
	if err != nil {
		return err
	}

	result := executor.New("tar", "--strip", "1", "-xzC", n.path, "-f", tarPath).Run()
	if result.Error != nil {
		return fmt.Errorf("failed to extract %s to %s: %s", tarPath, n.path, result.Error)
	}

	return nil
}
