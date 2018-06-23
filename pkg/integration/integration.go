package integration

import (
	"fmt"
	"os"
	"strings"

	color "github.com/logrusorgru/aurora"
)

// Print prints the integration code for the user's shell
func Print() {
	var currentShell = os.Getenv("SHELL")

	if currentShell == "" {
		currentShell = "bash"
		fmt.Fprintln(os.Stderr, color.Red("SHELL environment variable is empty"))
	}

	if strings.HasSuffix(currentShell, "bash") {
		fmt.Println(shellSource, bashSource)
	} else if strings.HasSuffix(currentShell, "zsh") {
		fmt.Println(shellSource, zshSource)
	} else {
		fmt.Fprintln(os.Stderr, color.Brown("Your shell is not supported"))
	}
}

// AddFinalizerCd declares a "cd" finalizer (change directory)
func AddFinalizerCd(path string) error {
	return addFinalizer("cd", path)
}

func addFinalizer(action, arg string) (err error) {
	content := fmt.Sprintf("%s:%s\n", action, arg)

	finalizerPath := os.Getenv("BUD_FINALIZER_FILE")

	if finalizerPath == "" {
		fmt.Println(color.Red("Shell integration error:"), "can't run a finalizer action:", color.Brown(content))
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
