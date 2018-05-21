package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type Executor struct {
	cmd              *exec.Cmd
	outputPrefix     string
	filterSubstrings []string
	outputWriter     io.Writer
}

// New returns an *Executor that will run the program with arguments
func New(program string, args ...string) *Executor {
	return &Executor{cmd: exec.Command(program, args...)}
}

// NewShell returns an *Executor that will run the command line in a shell
func NewShell(cmdline string) *Executor {
	return &Executor{cmd: exec.Command("sh", "-c", cmdline)}
}

// SetCwd changes the current working directory the command will be run in
func (e *Executor) SetCwd(cwd string) *Executor {
	e.cmd.Dir = cwd
	return e
}

// SetEnv changes the environment variables that will be used to run the command
func (e *Executor) SetEnv(env []string) *Executor {
	e.cmd.Env = env
	return e
}

// SetOutputPrefix sets a prefix for each line printed by the command
func (e *Executor) SetOutputPrefix(prefix string) *Executor {
	e.outputPrefix = prefix
	return e
}

// AddOutputFilter adds a substring to the list used to suppress lines printed by the command
func (e *Executor) AddOutputFilter(substring string) *Executor {
	e.filterSubstrings = append(e.filterSubstrings, substring)
	return e
}

func (e *Executor) getExitCode(err error) (int, error) {
	if err == nil {
		code := e.cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
		return code, nil
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		code := exitError.Sys().(syscall.WaitStatus).ExitStatus()
		return code, nil
	}

	// There was an error but not a ExitError, just return it with an invalid exit code
	return -1, err
}

func (e *Executor) printPipe(wg *sync.WaitGroup, pipe io.Reader) {
	defer wg.Done()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		if e.shouldSuppressLine(line) {
			continue
		}
		fmt.Fprintf(e.outputWriter, "%s%s\n", e.outputPrefix, line)
	}
}

func (e *Executor) shouldSuppressLine(line string) bool {
	for _, substring := range e.filterSubstrings {
		if strings.Contains(line, substring) {
			return true
		}
	}
	return false
}

// RunWithCode executes the command. Return the exit code and an error.
func (e *Executor) RunWithCode() (int, error) {
	if e.outputWriter == nil {
		e.outputWriter = os.Stdout
	}

	e.cmd.Stdin = nil
	stdout, err := e.cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}
	stderr, err := e.cmd.StderrPipe()
	if err != nil {
		return -1, err
	}

	err = e.cmd.Start()
	if err != nil {
		return -1, err
	}

	outputWait := new(sync.WaitGroup)
	outputWait.Add(2)
	go e.printPipe(outputWait, stdout)
	go e.printPipe(outputWait, stderr)
	outputWait.Wait()

	err = e.cmd.Wait()
	code, err := e.getExitCode(err)
	if err != nil {
		return code, fmt.Errorf("command failed with: %s", err)
	}
	return code, err
}

// Run executes the command. Return an error for non-zero exitcode.
func (e *Executor) Run() error {
	code, err := e.RunWithCode()
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("command failed with exit code %d", code)
	}
	return nil
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
	return strings.Trim(output, "\n"), code, err
}
