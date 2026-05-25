package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	api.Register("uv", "uv", parserUV).SetRequiredTask(pythonTaskName)
}

type uvConfig struct {
	Groups    []string
	Extras    []string
	AllGroups bool
	AllExtras bool
	Exact     bool
	Frozen    bool
}

func parserUV(config *api.TaskConfig, task *api.Task) error {
	uvCfg, err := parseUVConfig(config)
	if err != nil {
		return err
	}

	task.Info = "sync"
	parserUVInstall(task)
	parserUVSync(task, uvCfg)
	return nil
}

func parseUVConfig(config *api.TaskConfig) (uvConfig, error) {
	groups, err := config.GetListOfStringsPropertyDefault("groups", []string{})
	if err != nil {
		return uvConfig{}, err
	}
	extras, err := config.GetListOfStringsPropertyDefault("extras", []string{})
	if err != nil {
		return uvConfig{}, err
	}

	allGroups, err := getOptionalBool(config, "all_groups")
	if err != nil {
		return uvConfig{}, err
	}
	allExtras, err := getOptionalBool(config, "all_extras")
	if err != nil {
		return uvConfig{}, err
	}
	exact, err := getOptionalBool(config, "exact")
	if err != nil {
		return uvConfig{}, err
	}
	frozen, err := getOptionalBool(config, "frozen")
	if err != nil {
		return uvConfig{}, err
	}

	return uvConfig{
		Groups:    groups,
		Extras:    extras,
		AllGroups: allGroups,
		AllExtras: allExtras,
		Exact:     exact,
		Frozen:    frozen,
	}, nil
}

func getOptionalBool(config *api.TaskConfig, name string) (bool, error) {
	value, _, err := config.GetBooleanProperty(name)
	return value, err
}

func parserUVInstall(task *api.Task) {
	needed := func(ctx *context.Context) *api.ActionResult {
		version, err := findAutoEnvFeatureParam(ctx, pythonTaskName)
		if err != nil {
			return api.Failed("missing python feature: %s", err)
		}
		venv := helpers.NewVirtualenv(ctx.Cfg, helpers.VirtualenvName(ctx.Project, version))
		if !utils.PathExists(venv.Which("uv")) {
			return api.Needed("uv is not installed in the virtualenv")
		}
		return api.NotNeeded()
	}
	run := func(ctx *context.Context) error {
		pipArgs := []string{"install", "--require-virtualenv", "uv"}
		ctx.UI.TaskCommand("pip", pipArgs...)
		result := ctx.Executor.Run(executor.New("pip", pipArgs...))
		if result.Error != nil {
			return fmt.Errorf("failed to install uv: %w", result.Error)
		}
		return nil
	}
	task.AddActionBuilder("install uv", run).On(api.FuncCondition(needed))
}

func parserUVSync(task *api.Task, uvCfg uvConfig) {
	args := buildUVSyncArgs(uvCfg)
	run := func(ctx *context.Context) error {
		ctx.UI.TaskCommand("uv", args...)
		result := ctx.Executor.Run(executor.New("uv", args...))
		if result.Error != nil {
			return fmt.Errorf("uv sync failed: %w", result.Error)
		}
		return nil
	}
	task.AddActionBuilder("sync uv project", run).
		On(api.FuncCondition(uvProjectExists)).
		On(api.FileCondition("pyproject.toml")).
		On(api.FileCondition("uv.lock"))
}

func uvProjectExists(ctx *context.Context) *api.ActionResult {
	if fileExists(ctx, "pyproject.toml") {
		return api.NotNeeded()
	}
	return api.Failed("uv task requires pyproject.toml")
}

func buildUVSyncArgs(config uvConfig) []string {
	args := []string{"sync", "--active"}
	if !config.Exact {
		args = append(args, "--inexact")
	}
	for _, group := range config.Groups {
		args = append(args, "--group", group)
	}
	for _, extra := range config.Extras {
		args = append(args, "--extra", extra)
	}
	if config.AllGroups {
		args = append(args, "--all-groups")
	}
	if config.AllExtras {
		args = append(args, "--all-extras")
	}
	if config.Frozen {
		args = append(args, "--frozen")
	}
	return args
}
