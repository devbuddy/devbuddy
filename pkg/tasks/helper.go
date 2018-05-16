package tasks

import (
	"errors"
	"os"
	"path/filepath"

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

func runCommand(ctx *context, program string, args ...string) (int, error) {
	ctx.ui.TaskCommand(program, args...)
	return executor.New(program, args...).SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ()).Run()
}

func runShellSilent(ctx *context, cmdline string) (int, error) {
	return executor.NewShell(cmdline).SetCwd(ctx.proj.Path).SetEnv(ctx.env.Environ()).Run()
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
