package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/executor"
)

func asString(value interface{}) (string, error) {
	result, ok := value.(string)
	if ok {
		return result, nil
	}

	return "", fmt.Errorf("not a string: %T (%+v)", value, value)
}

func asListOfStrings(value interface{}) ([]string, error) {
	if v, ok := value.([]string); ok {
		return v, nil
	}

	elements, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a list of strings: type %T (%+v)", value, value)
	}

	listOfStrings := []string{}

	for _, element := range elements {
		str, ok := element.(string)
		if !ok {
			return nil, fmt.Errorf("not a list of strings: invalid element: type %T (%+v)", element, element)
		}
		listOfStrings = append(listOfStrings, str)
	}

	return listOfStrings, nil
}

func command(ctx *Context, program string, args ...string) executor.Executor {
	ctx.ui.TaskCommand(program, args...)
	return commandSilent(ctx, program, args...)
}

func sudoCommand(ctx *Context, program string, args ...string) executor.Executor {
	args = append([]string{program}, args...)
	program = "sudo"
	ctx.ui.TaskCommand(program, args...)
	return commandSilent(ctx, program, args...)
}

func commandSilent(ctx *Context, program string, args ...string) executor.Executor {
	return executor.New(program, args...).SetOutputPrefix("  ").SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ())
}

func shell(ctx *Context, cmdline string) executor.Executor {
	ctx.ui.TaskShell(cmdline)
	return shellSilent(ctx, cmdline)
}

func shellSilent(ctx *Context, cmdline string) executor.Executor {
	return executor.NewShell(cmdline).SetOutputPrefix("  ").SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ())
}

func fileExists(ctx *Context, path string) bool {
	if _, err := os.Stat(filepath.Join(ctx.proj.Path, path)); os.IsNotExist(err) {
		return false
	}
	return true
}

func fileModTime(ctx *Context, path string) (int64, error) {
	s, err := os.Stat(filepath.Join(ctx.proj.Path, path))
	if err != nil {
		return 0, err
	}
	return s.ModTime().UnixNano(), nil
}
