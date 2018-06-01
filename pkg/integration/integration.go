package integration

import (
	"fmt"
	"os"
	"strings"

	color "github.com/logrusorgru/aurora"
)

func Print() {
	var currentShell = os.Getenv("SHELL")

	if currentShell == "" {
		currentShell = "bash"
	}

	if strings.HasSuffix(currentShell, "bash") {
		fmt.Println(shellSource, bashSource)
	} else if strings.HasSuffix(currentShell, "zsh") {
		fmt.Println(shellSource, zshSource)
	} else {
		fmt.Fprintln(os.Stderr, color.Brown("Your shell is not supported"))
	}
}

func AddFinalizerCd(path string) error {
	return addFinalizer("cd", path)
}

func addFinalizer(action, arg string) (err error) {
	content := fmt.Sprintf("%s:%s\n", action, arg)

	finalizerPath := os.Getenv("DAD_FINALIZER_FILE")

	if finalizerPath == "" {
		fmt.Println(color.Red("Shell integration error:"), "can't run a finalizer action:", color.Brown(content))
		return nil
	}

	f, err := os.OpenFile(finalizerPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
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
