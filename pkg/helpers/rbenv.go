package helpers

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

// ReadRubyVersionFile reads a `.ruby-version` file and returns the version
// string with any leading `ruby-` engine prefix stripped. Returns os.ErrNotExist
// when the file is absent so callers can decide whether to treat that as fatal.
func ReadRubyVersionFile(projectPath string) (string, error) {
	data, err := os.ReadFile(filepath.Join(projectPath, ".ruby-version"))
	if err != nil {
		return "", err
	}
	version := strings.TrimSpace(string(data))
	version = strings.TrimPrefix(version, "ruby-")
	if version == "" {
		return "", fmt.Errorf(".ruby-version is empty")
	}
	return version, nil
}

type RbEnv struct {
	ctx     *context.Context
	command string
	root    string
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
	command := "rbenv"
	result := ctx.Executor.CaptureAndTrim(executor.New(command, "root"))
	if result.Error != nil && utils.PathExists(sourceRbEnvCommand()) {
		command = sourceRbEnvCommand()
		result = ctx.Executor.CaptureAndTrim(executor.New(command, "root"))
	}
	if result.Error != nil {
		return nil, fmt.Errorf("Command '%s root' failed: %w", command, result.Error)
	}
	return &RbEnv{ctx: ctx, command: command, root: result.Output}, nil
}

func sourceRbEnvCommand() string {
	return path.Join(RbEnvRoot(), "bin", "rbenv")
}

func (r *RbEnv) Command() string {
	return r.command
}

func (r *RbEnv) VersionInstalled(version string) (bool, error) {
	versions, err := r.listVersions()
	if err != nil {
		return false, err
	}
	return slices.Contains(versions, version), nil
}

func (r *RbEnv) listVersions() ([]string, error) {
	result := r.ctx.Executor.Capture(executor.New(r.command, "versions", "--bare", "--skip-aliases"))
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
