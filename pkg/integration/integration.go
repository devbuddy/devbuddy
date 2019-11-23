package integration

import (
	"fmt"
	"os"
	"strings"

	ps "github.com/mitchellh/go-ps"

	"github.com/devbuddy/devbuddy/pkg/termui"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

// Print prints the integration code for the user's shell
func Print() {
	shell, err := detectShell()
	if err != nil {
		termui.HookShellDetectionError(err)
	}

	switch shell {
	case "bash":
		fmt.Println(shellSource, bashSource)
	case "zsh":
		fmt.Println(shellSource, zshSource)
	}
}

func detectShell() (string, error) {
	proc, err := ps.FindProcess(os.Getppid())
	if err != nil {
		return "", fmt.Errorf("failed to get parent process info: %s", err)
	}
	parentProcessPath := proc.Executable()

	switch {
	case strings.HasSuffix(parentProcessPath, "bash"):
		return "bash", nil
	case strings.HasSuffix(parentProcessPath, "zsh"):
		return "zsh", nil
	}

	return "", fmt.Errorf("parent process is not a supported shell: %s", parentProcessPath)
}

// AddFinalizerCd declares a "cd" finalizer (change directory)
func AddFinalizerCd(path string) error {
	return addFinalizer("cd", path)
}

func formatError(message string) error {
	return fmt.Errorf(`there is something wrong this the shell integration:

    %s

This usually means that DevBuddy is not setup properly.
Please follow the setup steps: https://github.com/devbuddy/devbuddy/tree/master#setup

If DevBuddy is already setup, then please open an issue on https://github.com/devbuddy/devbuddy/issues/new?labels=bug
You can use "bud --report-issue" to do that.
`, message)
}

func addFinalizer(action, arg string) (err error) {
	content := fmt.Sprintf("%s:%s\n", action, arg)

	finalizerPath := os.Getenv("BUD_FINALIZER_FILE")

	if finalizerPath == "" {
		return formatError("the BUD_FINALIZER_FILE environment variable is missing or empty")
	}

	return utils.AppendOnlyFile(finalizerPath, []byte(content))
}
