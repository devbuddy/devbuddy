package tasks

import (
	"errors"

	"github.com/pior/dad/pkg/executor"
)

func asString(value interface{}) (string, error) {
	result, ok := value.(string)
	if ok {
		return result, nil
	}

	_, ok = value.(bool)
	if ok {
		return "", errors.New("found a boolean, not a string")
	}

	return "", errors.New("not a string")
}

func runCommand(ctx *Context, program string, args ...string) (int, error) {
	ctx.ui.TaskCommand(program, args...)
	return executor.New(program, args...).SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ()).Run()
}

func runShellSilent(ctx *Context, cmdline string) (int, error) {
	return executor.NewShell(cmdline).SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ()).Run()
}
