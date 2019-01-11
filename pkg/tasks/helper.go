package tasks

import (
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func command(ctx *taskapi.Context, program string, args ...string) *executor.Executor {
	ctx.ui.TaskCommand(program, args...)
	return commandSilent(ctx, program, args...)
}

func sudoCommand(ctx *taskapi.Context, program string, args ...string) *executor.Executor {
	args = append([]string{program}, args...)
	program = "sudo"
	ctx.ui.TaskCommand(program, args...)
	return commandSilent(ctx, program, args...)
}

func commandSilent(ctx *taskapi.Context, program string, args ...string) *executor.Executor {
	return executor.New(program, args...).SetOutputPrefix("  ").SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ())
}

func shell(ctx *taskapi.Context, cmdline string) *executor.Executor {
	ctx.ui.TaskShell(cmdline)
	return shellSilent(ctx, cmdline)
}

func shellSilent(ctx *taskapi.Context, cmdline string) *executor.Executor {
	return executor.NewShell(cmdline).SetOutputPrefix("  ").SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ())
}

func fileExists(ctx *taskapi.Context, path string) bool {
	if _, err := os.Stat(filepath.Join(ctx.proj.Path, path)); os.IsNotExist(err) {
		return false
	}
	return true
}

func fileModTime(ctx *taskapi.Context, path string) (int64, error) {
	s, err := os.Stat(filepath.Join(ctx.proj.Path, path))
	if err != nil {
		return 0, err
	}
	return s.ModTime().UnixNano(), nil
}
