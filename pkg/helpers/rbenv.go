package helpers

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
)

type RbEnv struct {
	ctx  *context.Context
	root string
}

// RbEnvRoot returns the rbenv root directory using the same resolution as
// rbenv itself: $RBENV_ROOT if set, otherwise $HOME/.rbenv. This avoids
// shelling out to `rbenv root` for hot paths like shell-hook activation.
func RbEnvRoot() string {
	if root := os.Getenv("RBENV_ROOT"); root != "" {
		return root
	}
	return path.Join(os.Getenv("HOME"), ".rbenv")
}

func NewRbEnv(ctx *context.Context) (*RbEnv, error) {
	result := ctx.Executor.CaptureAndTrim(executor.New("rbenv", "root"))
	if result.Error != nil {
		return nil, fmt.Errorf("Command 'rbenv root' failed: %w", result.Error)
	}
	return &RbEnv{ctx: ctx, root: result.Output}, nil
}

func (r *RbEnv) VersionInstalled(version string) (bool, error) {
	versions, err := r.listVersions()
	if err != nil {
		return false, err
	}
	return slices.Contains(versions, version), nil
}

func (r *RbEnv) listVersions() ([]string, error) {
	result := r.ctx.Executor.Capture(executor.New("rbenv", "versions", "--bare", "--skip-aliases"))
	if result.Error != nil {
		return nil, fmt.Errorf("failed to run rbenv versions: %w", result.Error)
	}
	return strings.Split(strings.TrimSpace(result.Output), "\n"), nil
}

func (r *RbEnv) VersionPath(version string) string {
	return path.Join(r.root, "versions", version)
}

func (r *RbEnv) BinPath(version string) string {
	return path.Join(r.VersionPath(version), "bin")
}

func (r *RbEnv) Which(version string, command string) string {
	return path.Join(r.BinPath(version), command)
}
