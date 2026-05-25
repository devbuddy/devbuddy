package tasks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/helpers/osidentity"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

const rubyTaskName = "ruby"

func init() {
	api.Register(rubyTaskName, "Ruby", parserRuby)
}

func parserRuby(config *api.TaskConfig, task *api.Task) error {
	version, err := resolveRubyVersion(config)
	if err != nil {
		return err
	}
	task.Info = version

	parserRubyInstallRbenv(task)
	parserRubyInstallRubyVersion(task, version)
	parserRubyBundleInstall(task, version)
	return nil
}

// resolveRubyVersion returns the requested Ruby version. It accepts the version
// from the dev.yml payload (string form or `version:` property) and falls back
// to a `.ruby-version` file at the project root when neither is provided.
func resolveRubyVersion(config *api.TaskConfig) (string, error) {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err == nil {
		return version, nil
	}

	if config.ProjectPath != "" {
		v, readErr := helpers.ReadRubyVersionFile(config.ProjectPath)
		if readErr == nil {
			return v, nil
		}
		if !errors.Is(readErr, os.ErrNotExist) {
			return "", fmt.Errorf("reading .ruby-version: %w", readErr)
		}
	}
	return "", err
}

func parserRubyInstallRbenv(task *api.Task) {
	needed := func(ctx *context.Context) *api.ActionResult {
		rbEnv, err := helpers.NewRbEnv(ctx)
		if err != nil {
			return api.Needed("rbenv is not installed: %s", err)
		}
		hasRubyBuild, err := rbenvHasInstallCommand(ctx, rbEnv)
		if err != nil {
			return api.Failed("failed to inspect rbenv commands: %s", err)
		}
		if !hasRubyBuild {
			return api.Needed("ruby-build is not installed")
		}
		return api.NotNeeded()
	}
	run := func(ctx *context.Context) error {
		return installRbenv(ctx, osidentity.Detect())
	}
	task.AddActionBuilder("install rbenv", run).On(api.FuncCondition(needed))
}

func rbenvHasInstallCommand(ctx *context.Context, rbEnv *helpers.RbEnv) (bool, error) {
	result := ctx.Executor.Capture(executor.New(rbEnv.Command(), "commands"))
	if result.Error != nil {
		return false, result.Error
	}
	for _, command := range strings.Split(result.Output, "\n") {
		if strings.TrimSpace(command) == "install" {
			return true, nil
		}
	}
	return false, nil
}

func installRbenv(ctx *context.Context, osIdent *osidentity.Identity) error {
	switch {
	case osIdent.IsMacOS():
		return installRbenvWithHomebrew(ctx)
	case osIdent.IsDebianLike():
		return installRbenvWithApt(ctx)
	default:
		return installRbenvFromGit(ctx)
	}
}

func installRbenvWithHomebrew(ctx *context.Context) error {
	result := ctx.RunTaskCommand(executor.New("brew", "install", "rbenv", "ruby-build").AddEnvVar("HOMEBREW_NO_AUTO_UPDATE", "1"))
	if result.Error != nil {
		return fmt.Errorf("failed to install rbenv: %w", result.Error)
	}
	return nil
}

func installRbenvWithApt(ctx *context.Context) error {
	result := ctx.RunTaskCommand(executor.New("sudo", "apt-get", "update"))
	if result.Error != nil {
		return fmt.Errorf("failed to run apt-get update: %w", result.Error)
	}

	result = ctx.RunTaskCommand(executor.New("sudo", "apt-get", "install", "--no-install-recommends", "-y", "rbenv", "ruby-build"))
	if result.Error != nil {
		return fmt.Errorf("failed to install rbenv: %w", result.Error)
	}
	return nil
}

func installRbenvFromGit(ctx *context.Context) error {
	root := helpers.RbEnvRoot()
	if !utils.PathExists(root) {
		result := ctx.RunTaskCommand(executor.New("git", "clone", "https://github.com/rbenv/rbenv.git", root))
		if result.Error != nil {
			return fmt.Errorf("failed to clone rbenv: %w", result.Error)
		}
	}

	rubyBuildPath := filepath.Join(root, "plugins", "ruby-build")
	if !utils.PathExists(rubyBuildPath) {
		result := ctx.RunTaskCommand(executor.New("git", "clone", "https://github.com/rbenv/ruby-build.git", rubyBuildPath))
		if result.Error != nil {
			return fmt.Errorf("failed to clone ruby-build: %w", result.Error)
		}
	}
	return nil
}

func parserRubyInstallRubyVersion(task *api.Task, version string) {
	needed := func(ctx *context.Context) *api.ActionResult {
		rbEnv, err := helpers.NewRbEnv(ctx)
		if err != nil {
			return api.Failed("cannot use rbenv: %s", err)
		}
		installed, err := rbEnv.VersionInstalled(version)
		if err != nil {
			return api.Failed("failed to check if ruby version is installed: %s", err)
		}
		if !installed {
			return api.Needed("ruby version is not installed")
		}
		return api.NotNeeded()
	}
	run := func(ctx *context.Context) error {
		if err := helpers.EnsureXcodeCommandLineTools(ctx); err != nil {
			return err
		}
		rbEnv, err := helpers.NewRbEnv(ctx)
		if err != nil {
			return fmt.Errorf("cannot use rbenv: %w", err)
		}
		result := ctx.RunTaskCommand(executor.New(rbEnv.Command(), "install", version))
		if result.Error != nil {
			return fmt.Errorf("failed to install the required ruby version: %w", result.Error)
		}
		return nil
	}
	task.AddActionBuilder("install Ruby version with rbenv", run).
		On(api.FuncCondition(needed)).
		SetFeature("ruby", version)
}

func parserRubyBundleInstall(task *api.Task, version string) {
	run := func(ctx *context.Context) error {
		rbEnv, err := helpers.NewRbEnv(ctx)
		if err != nil {
			return err
		}
		bundle := rbEnv.Which(version, "bundle")
		result := ctx.RunTaskCommand(executor.New(bundle, "config", "set", "--local", "path", "vendor/bundle"))
		if result.Error != nil {
			return fmt.Errorf("bundle config failed: %w", result.Error)
		}
		result = ctx.RunTaskCommand(executor.New(bundle, "install"))
		if result.Error != nil {
			return fmt.Errorf("bundle install failed: %w", result.Error)
		}
		return nil
	}
	// Either file changing should trigger a re-run: Gemfile when deps are
	// added/removed, Gemfile.lock when versions are bumped (e.g. `bundle update`).
	task.AddActionBuilder("install gems with bundler", run).
		On(api.FileCondition("Gemfile")).
		On(api.FileCondition("Gemfile.lock"))
}
