package executor

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Executor struct {
	cmd *exec.Cmd
}

// New returns an *Executor that will run the program with arguments
func New(program string, args ...string) *Executor {
	return &Executor{cmd: exec.Command(program, args...)}
}

// NewShell returns an *Executor that will run the command line in a shell
func NewShell(cmdline string) *Executor {
	return &Executor{cmd: exec.Command("sh", "-c", cmdline)}
}

// SetCwd change the current working directory the command will be run in
func (e *Executor) SetCwd(cwd string) *Executor {
	e.cmd.Dir = cwd
	return e
}

func (e *Executor) getExitCode(err error) (int, error) {
	if err == nil {
		code := e.cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
		return code, nil
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		code := exitError.Sys().(syscall.WaitStatus).ExitStatus()
		return code, err
	}

	// There was an error but not a ExitError, just return it with an invalid exit code
	return -1, err
}

// Run executes the command and returns the exit code
func (e *Executor) Run() (int, error) {
	e.cmd.Stdin = os.Stdin
	e.cmd.Stdout = os.Stdout
	e.cmd.Stderr = os.Stderr

	err := e.cmd.Run()
	code, err := e.getExitCode(err)
	return code, err
}

// Capture executes the command and return the output and the exit code
func (e *Executor) Capture() (string, int, error) {
	output, err := e.cmd.Output()
	code, err := e.getExitCode(err)
	return string(output), code, err
}

// CaptureAndTrim calls Capture() and trim the blank lines
func (e *Executor) CaptureAndTrim() (string, int, error) {
	output, code, err := e.Capture()
	return strings.Trim(string(output), "\n"), code, err
}
