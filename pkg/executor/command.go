package executor

import (
	"io"

	"github.com/devbuddy/devbuddy/pkg/env"
)

// Command describes what to execute (pure data, inspectable by tests).
type Command struct {
	Program       string
	Args          []string
	Shell         bool // if true, Program is a shell cmdline run via sh -c
	Cwd           string
	Env           []string
	Passthrough   bool
	OutputPrefix  string
	OutputFilters []string
	OutputWriter  io.Writer // nil defaults to os.Stdout
}

// New returns a Command that will run the program with arguments.
func New(program string, args ...string) *Command {
	return &Command{Program: program, Args: args}
}

// NewShell returns a Command that will run the command line in a shell.
func NewShell(cmdline string) *Command {
	return &Command{Program: cmdline, Shell: true}
}

// AddEnvVar sets a single variable in the command's environment.
func (c *Command) AddEnvVar(name, value string) *Command {
	environ := env.New(c.Env)
	environ.Set(name, value)
	c.Env = environ.Environ()
	return c
}

// AddOutputFilter adds a substring to the list used to suppress output lines.
func (c *Command) AddOutputFilter(substring string) *Command {
	c.OutputFilters = append(c.OutputFilters, substring)
	return c
}
