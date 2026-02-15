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

// Runner is the mock point for command execution.
type Runner interface {
	Run(cmd *Command) *Result
	Capture(cmd *Command) *Result
}

// Executor runs commands using a Runner, applying defaults (Cwd, Env, OutputPrefix).
type Executor struct {
	Runner       Runner
	Cwd          string
	Env          *env.Env // live reference, evaluated at run time
	OutputPrefix string
}

// NewExecutor returns an Executor backed by a real os/exec Runner.
func NewExecutor() *Executor {
	return &Executor{Runner: &osRunner{}}
}

func (e *Executor) applyDefaults(cmd *Command) {
	if cmd.Cwd == "" && e.Cwd != "" {
		cmd.Cwd = e.Cwd
	}
	if e.Env != nil {
		if cmd.Env == nil {
			cmd.Env = e.Env.Environ()
		} else {
			cmd.Env = env.MergeEnviron(e.Env.Environ(), cmd.Env)
		}
	}
	if cmd.OutputPrefix == "" && e.OutputPrefix != "" {
		cmd.OutputPrefix = e.OutputPrefix
	}
}

// Run executes the command with visible output.
func (e *Executor) Run(cmd *Command) *Result {
	e.applyDefaults(cmd)
	return e.Runner.Run(cmd)
}

// Capture executes the command and captures its stdout.
func (e *Executor) Capture(cmd *Command) *Result {
	e.applyDefaults(cmd)
	return e.Runner.Capture(cmd)
}

// CaptureAndTrim captures and trims trailing newlines from the output.
func (e *Executor) CaptureAndTrim(cmd *Command) *Result {
	r := e.Capture(cmd)
	r.Output = strings.Trim(r.Output, "\n")
	return r
}

// osRunner implements Runner via os/exec.
type osRunner struct{}

func buildExecCmd(cmd *Command) *exec.Cmd {
	var c *exec.Cmd
	if cmd.Shell {
		c = exec.Command("sh", "-c", cmd.Program)
	} else {
		c = exec.Command(cmd.Program, cmd.Args...)
	}
	if cmd.Cwd != "" {
		c.Dir = cmd.Cwd
	}
	if cmd.Env != nil {
		c.Env = cmd.Env
	}
	return c
}

func (r *osRunner) Run(cmd *Command) *Result {
	c := buildExecCmd(cmd)

	outputWriter := cmd.OutputWriter
	if outputWriter == nil {
		outputWriter = os.Stdout
	}

	var err error
	switch {
	case cmd.Passthrough:
		err = runPassthrough(c)
	default:
		err = runWithOutputFilter(cmd, c, outputWriter)
	}
	return buildResult("", err)
}

func (r *osRunner) Capture(cmd *Command) *Result {
	c := buildExecCmd(cmd)
	output, err := c.Output()
	return buildResult(string(output), err)
}

func runWithOutputFilter(cmd *Command, c *exec.Cmd, outputWriter io.Writer) error {
	c.Stdin = nil // default to /dev/null

	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	if err = c.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go printPipe(&wg, stdout, cmd.OutputPrefix, cmd.OutputFilters, outputWriter)
	go printPipe(&wg, stderr, cmd.OutputPrefix, cmd.OutputFilters, outputWriter)
	wg.Wait()

	return c.Wait()
}

func printPipe(wg *sync.WaitGroup, pipe io.Reader, prefix string, filters []string, w io.Writer) {
	defer wg.Done()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		if shouldSuppressLine(line, filters) {
			continue
		}
		termui.Fprintf(w, "%s%s\n", prefix, line)
	}
}

func shouldSuppressLine(line string, filters []string) bool {
	for _, substring := range filters {
		if strings.Contains(line, substring) {
			return true
		}
	}
	return false
}
