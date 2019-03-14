package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

// RegisterTasks is a hack to force the execution of the task registration (in the init functions)
func RegisterTasks() {}

func command(ctx *context.Context, program string, args ...string) executor.Executor {
	ctx.UI.TaskCommand(program, args...)
	return commandSilent(ctx, program, args...)
}

func sudoCommand(ctx *context.Context, program string, args ...string) executor.Executor {
	args = append([]string{program}, args...)
	program = "sudo"
	ctx.UI.TaskCommand(program, args...)
	return commandSilent(ctx, program, args...)
}

func commandSilent(ctx *context.Context, program string, args ...string) executor.Executor {
	return executor.New(program, args...).SetOutputPrefix("  ").SetCwd(ctx.Project.Path).SetEnv(ctx.Env.Environ())
}

func shell(ctx *context.Context, cmdline string) executor.Executor {
	ctx.UI.TaskShell(cmdline)
	return shellSilent(ctx, cmdline)
}

func shellSilent(ctx *context.Context, cmdline string) executor.Executor {
	return executor.NewShell(cmdline).SetOutputPrefix("  ").SetCwd(ctx.Project.Path).SetEnv(ctx.Env.Environ())
}

func fileExists(ctx *context.Context, path string) bool {
	if _, err := os.Stat(filepath.Join(ctx.Project.Path, path)); os.IsNotExist(err) {
		return false
	}
	return true
}

func fileModTime(ctx *context.Context, path string) (int64, error) {
	s, err := os.Stat(filepath.Join(ctx.Project.Path, path))
	if err != nil {
		return 0, err
	}
	return s.ModTime().UnixNano(), nil
}

func findAutoEnvFeatureParam(ctx *context.Context, name string) (string, error) {
	taskList, err := taskapi.GetTasksFromProject(ctx.Project)
	if err != nil {
		return "", err
	}
	feature := taskapi.GetFeaturesFromTasks(taskList).Get(name)
	if feature == nil {
		return "", fmt.Errorf("no autoenv feature with name %s", name)
	}
	return feature.Param, nil
}
