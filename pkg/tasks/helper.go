package tasks

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/executor"
)

func asString(value interface{}) (string, error) {
	result, ok := value.(string)
	if ok {
		return result, nil
	}

	_, ok = value.(bool)
	if ok {
		return "", errors.New("not a string")
	}

	return "", errors.New("not a string")
}

func command(ctx *context, program string, args ...string) *executor.Executor {
	ctx.ui.TaskCommand(program, args...)
	return commandSilent(ctx, program, args...)
}

func commandSilent(ctx *context, program string, args ...string) *executor.Executor {
	return executor.New(program, args...).SetOutputPrefix("  ").SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ())
}

func shell(ctx *context, cmdline string) *executor.Executor {
	ctx.ui.TaskShell(cmdline)
	return shellSilent(ctx, cmdline)
}

func shellSilent(ctx *context, cmdline string) *executor.Executor {
	return executor.NewShell(cmdline).SetOutputPrefix("  ").SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ())
}

func fileExists(ctx *context, path string) bool {
	if _, err := os.Stat(filepath.Join(ctx.proj.Path, path)); os.IsNotExist(err) {
		return false
	}
	return true
}

func fileModTime(ctx *context, path string) (int64, error) {
	s, err := os.Stat(filepath.Join(ctx.proj.Path, path))
	if err != nil {
		return 0, err
	}
	return s.ModTime().UnixNano(), nil
}
