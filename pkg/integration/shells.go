package integration

import (
	"fmt"
	"os"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

type ShellIdentity string

var (
	BASH ShellIdentity = "bash"
	ZSH  ShellIdentity = "zsh"
)

func DetectShell() (ShellIdentity, error) {
	proc, err := ps.FindProcess(os.Getppid())
	if err != nil {
		return "", fmt.Errorf("failed to get parent process info: %w", err)
	}
	parentProcessPath := proc.Executable()

	switch {
	case strings.HasSuffix(parentProcessPath, "bash"):
		return BASH, nil
	case strings.HasSuffix(parentProcessPath, "zsh"):
		return ZSH, nil
	}

	return "", fmt.Errorf("parent process is not a supported shell: %s", parentProcessPath)
}
