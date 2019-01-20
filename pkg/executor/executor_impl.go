package executor

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

type executorImpl struct {
	cmd              *exec.Cmd
	outputPrefix     string
	filterSubstrings []string
	outputWriter     io.Writer
}

// New returns an Executor that will run the program with arguments
func New(program string, args ...string) Executor {
	return &executorImpl{cmd: exec.Command(program, args...)}
}

// NewShell returns an Executor that will run the command line in a shell
func NewShell(cmdline string) Executor {
	return &executorImpl{cmd: exec.Command("sh", "-c", cmdline)}
}

// SetCwd changes the current working directory the command will be run in
func (e *executorImpl) SetCwd(cwd string) Executor {
	e.cmd.Dir = cwd
	return e
}

// SetEnv changes the environment variables that will be used to run the command
func (e *executorImpl) SetEnv(env []string) Executor {
	e.cmd.Env = env
	return e
}

// SetEnvVar sets a single variable in the environment that will be used to run the command
func (e *executorImpl) SetEnvVar(name, value string) Executor {
	env := env.New(e.cmd.Env)
	env.Set(name, value)
	e.cmd.Env = env.Environ()
	return e
}

// SetOutputPrefix sets a prefix for each line printed by the command
func (e *executorImpl) SetOutputPrefix(prefix string) Executor {
	e.outputPrefix = prefix
	return e
}

// AddOutputFilter adds a substring to the list used to suppress lines printed by the command
func (e *executorImpl) AddOutputFilter(substring string) Executor {
	e.filterSubstrings = append(e.filterSubstrings, substring)
	return e
}

func (e *executorImpl) printPipe(wg *sync.WaitGroup, pipe io.Reader) {
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

func (e *executorImpl) shouldSuppressLine(line string) bool {
	for _, substring := range e.filterSubstrings {
		if strings.Contains(line, substring) {
			return true
		}
	}
	return false
}

func (e *executorImpl) runWithOutputFilter() error {
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

// Run executes the command and returns a Result
func (e *executorImpl) Run() *Result {
	err := e.runWithOutputFilter()
	return buildResult("", err)
}

// Capture executes the command and return a Result
func (e *executorImpl) Capture() *Result {
	output, err := e.cmd.Output()
	return buildResult(string(output), err)
}

// CaptureAndTrim calls Capture() and trim the blank lines
func (e *executorImpl) CaptureAndTrim() *Result {
	output, err := e.cmd.Output()
	return buildResult(strings.Trim(string(output), "\n"), err)
}
