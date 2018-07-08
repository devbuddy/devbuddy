package integration

import (
	"fmt"
	"os"
	"strings"

	ps "github.com/mitchellh/go-ps"

	"github.com/devbuddy/devbuddy/pkg/termui"
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

func addFinalizer(action, arg string) (err error) {
	content := fmt.Sprintf("%s:%s\n", action, arg)

	finalizerPath := os.Getenv("BUD_FINALIZER_FILE")

	if finalizerPath == "" {
		termui.HookIntegrationError("can't run a finalizer action: " + content)
		return nil
	}

	return writeFile(finalizerPath, content)
}

func writeFile(path string, content string) (err error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = f.WriteString(content)
	if err != nil {
		return
	}

	return
}
