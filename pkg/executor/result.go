package executor

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

// Result represents the result of a command execution
type Result struct {
	Code        int    // code returned by the process
	Error       error  // error about the process launch and exit
	LaunchError error  // error about the process launch
	Output      string // command output if captured, otherwise empty
}

func buildResult(output string, err error) *Result {
	code := 0

	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			code = exitError.Sys().(syscall.WaitStatus).ExitStatus()
			err = nil
		} else {
			err = fmt.Errorf("command failed with: %w", err)
		}
	}

	errForCode := err
	if code > 0 {
		errForCode = fmt.Errorf("command failed with exit code %d", code)
	}

	return &Result{
		Error:       errForCode,
		LaunchError: err,
		Code:        code,
		Output:      output,
	}
}
