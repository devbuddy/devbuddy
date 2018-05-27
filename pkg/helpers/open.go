package helpers

import (
	"os/exec"
	"runtime"
)

// Open a file or URL with the default application, return immediately
func Open(location string) error {
	openCommand := "xdg-open"
	if runtime.GOOS == "darwin" {
		openCommand = "open"
	}

	return exec.Command(openCommand, location).Start()
}
