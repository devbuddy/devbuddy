package helpers

import (
	"fmt"
	"path"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/executor"
)

type PyEnv struct {
	root string
}

func NewPyEnv() (*PyEnv, error) {
	root, err := executor.New("pyenv", "root").CaptureAndTrim()
	if err != nil {
		return nil, fmt.Errorf("Command 'pyenv root' failed: %s", err)
	}
	v := PyEnv{root: root}
	return &v, nil
}

func (p *PyEnv) VersionInstalled(version string) (installed bool, err error) {
	versions, err := p.listVersions()
	if err != nil {
		return
	}

	for _, v := range versions {
		if v == version {
			return true, nil
		}
	}
	return
}

func (p *PyEnv) listVersions() ([]string, error) {
	output, err := executor.New("pyenv", "versions", "--bare", "--skip-aliases").Capture()
	if err != nil {
		return nil, fmt.Errorf("failed to run pyenv versions: %s", err)
	}

	versions := strings.Split(strings.TrimSpace(output), "\n")
	return versions, nil
}

func (p *PyEnv) Which(version string, command string) string {
	return path.Join(p.root, "versions", version, "bin", command)
}
