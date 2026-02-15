package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

// RegisterTasks is a hack to force the execution of the task registration (in the init functions)
func RegisterTasks() {}

func command(ctx *context.Context, program string, args ...string) *executor.Command {
	ctx.UI.TaskCommand(program, args...)
	return commandSilent(ctx, program, args...)
}

func sudoCommand(ctx *context.Context, program string, args ...string) *executor.Command {
	return command(ctx, "sudo", append([]string{program}, args...)...)
}

func commandSilent(ctx *context.Context, program string, args ...string) *executor.Command {
	return &executor.Command{
		Program:      program,
		Args:         args,
		OutputPrefix: "  ",
		Cwd:          ctx.Project.Path,
		Env:          ctx.Env.Environ(),
	}
}

func shell(ctx *context.Context, cmdline string) *executor.Command {
	ctx.UI.TaskShell(cmdline)
	return shellSilent(ctx, cmdline)
}

func shellSilent(ctx *context.Context, cmdline string) *executor.Command {
	return &executor.Command{
		Program:      cmdline,
		Shell:        true,
		OutputPrefix: "  ",
		Cwd:          ctx.Project.Path,
		Env:          ctx.Env.Environ(),
	}
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
	taskList, err := api.GetTasksFromProject(ctx.Project)
	if err != nil {
		return "", err
	}
	feature := api.GetFeaturesFromTasks(taskList).Get(name)
	if feature == nil {
		return "", fmt.Errorf("no autoenv feature with name %s", name)
	}
	return feature.Param, nil
}
