package helpers

import (
	"fmt"
	"path"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
)

type PyEnv struct {
	ctx  *context.Context
	root string
}

func NewPyEnv(ctx *context.Context) (*PyEnv, error) {
	result := ctx.Executor.CaptureAndTrim(executor.New("pyenv", "root"))
	if result.Error != nil {
		return nil, fmt.Errorf("Command 'pyenv root' failed: %w", result.Error)
	}
	v := PyEnv{ctx: ctx, root: result.Output}
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
	result := p.ctx.Executor.Capture(executor.New("pyenv", "versions", "--bare", "--skip-aliases"))
	if result.Error != nil {
		return nil, fmt.Errorf("failed to run pyenv versions: %w", result.Error)
	}

	versions := strings.Split(strings.TrimSpace(result.Output), "\n")
	return versions, nil
}

func (p *PyEnv) Which(version string, command string) string {
	return path.Join(p.root, "versions", version, "bin", command)
}
