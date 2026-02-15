package executor

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/devbuddy/devbuddy/pkg/termui"
)

// Executor runs commands.
type Executor interface {
	Run(cmd *Command) *Result
	Capture(cmd *Command) *Result
}

type realExecutor struct{}

// NewExecutor returns a real Executor that runs commands via os/exec.
func NewExecutor() Executor { return &realExecutor{} }

// Package-level convenience functions using the real executor.

func Run(cmd *Command) *Result { return (&realExecutor{}).Run(cmd) }

func Capture(cmd *Command) *Result { return (&realExecutor{}).Capture(cmd) }

func CaptureAndTrim(cmd *Command) *Result {
	r := Capture(cmd)
	r.Output = strings.Trim(r.Output, "\n")
	return r
}

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

func (e *realExecutor) Run(cmd *Command) *Result {
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

func (e *realExecutor) Capture(cmd *Command) *Result {
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
