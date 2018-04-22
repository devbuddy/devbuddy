package helpers

import (
	"fmt"
	"path"
	"strings"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
)

type PyEnv struct {
	root string
}

func NewPyEnv(cfg *config.Config) (*PyEnv, error) {
	root, code, err := executor.New("pyenv", "root").CaptureAndTrim()
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, fmt.Errorf("Command 'pyenv root' failed with code %d", code)
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
	output, code, err := executor.Capture("pyenv", "versions", "--bare", "--skip-aliases")
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, fmt.Errorf("failed to run pyenv versions. exit code: %d", code)
	}

	versions := strings.Split(strings.TrimSpace(output), "\n")
	return versions, nil
}

func (p *PyEnv) Which(version string, command string) string {
	return path.Join(p.root, "versions", version, "bin", command)
}
