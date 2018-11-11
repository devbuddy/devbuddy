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

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

// Executor prepares and run a command execution
type Executor struct {
	cmd              *exec.Cmd
	outputPrefix     string
	filterSubstrings []string
	outputWriter     io.Writer
}

// Result represents the result of a command execution
type Result struct {
	Code   int
	Error  error
	Output string
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

// SetEnvVar sets a single variable in the environment that will be used to run the command
func (e *Executor) SetEnvVar(name, value string) *Executor {
	env := env.New(e.cmd.Env)
	env.Set(name, value)
	e.cmd.Env = env.Environ()
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
		termui.Fprintf(e.outputWriter, "%s%s\n", e.outputPrefix, line)
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

func (e *Executor) runWithOutputFilter() error {
	if e.outputWriter == nil {
		e.outputWriter = os.Stdout
	}

	e.cmd.Stdin = nil
	stdout, err := e.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := e.cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = e.cmd.Start()
	if err != nil {
		return err
	}

	outputWait := new(sync.WaitGroup)
	outputWait.Add(2)
	go e.printPipe(outputWait, stdout)
	go e.printPipe(outputWait, stderr)
	outputWait.Wait()

	return e.cmd.Wait()
}

func (e *Executor) buildResult(output string, err error) *Result {
	code, err := e.getExitCode(err)
	if err != nil {
		err = fmt.Errorf("command failed with: %s", err)
	} else if code != 0 {
		err = fmt.Errorf("command failed with exit code %d", code)
	}
	return &Result{
		Error:  err,
		Code:   code,
		Output: output,
	}
}

// Run executes the command. Returns a Result. Code is -1 if the command failed to start
func (e *Executor) Run() *Result {
	err := e.runWithOutputFilter()
	return e.buildResult("", err)
}

// Capture executes the command and return a Result
func (e *Executor) Capture() *Result {
	output, err := e.cmd.Output()
	return e.buildResult(string(output), err)
}

// CaptureAndTrim calls Capture() and trim the blank lines
func (e *Executor) CaptureAndTrim() *Result {
	output, err := e.cmd.Output()
	return e.buildResult(strings.Trim(string(output), "\n"), err)
}
