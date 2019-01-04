package executor

import "github.com/devbuddy/devbuddy/pkg/env"

// ExecutorMock implement base Executor api for command line
type ExecutorMock struct {
	program          string
	args             []string
	outputPrefix     string
	workingDir       string
	env              []string
	filterSubstrings []string
}

func NewMock(program string, args ...string) *ExecutorMock {
	return &ExecutorMock{program: program, args: args}
}

func (e *ExecutorMock) SetCwd(cwd string) *ExecutorMock {
	e.workingDir = cwd
	return e
}

func (e *ExecutorMock) SetEnv(env []string) *ExecutorMock {
	e.env = env
	return e
}

func (e *ExecutorMock) SetEnvVar(name, value string) *ExecutorMock {
	env := env.New(e.env)
	env.Set(name, value)
	e.env = env.Environ()
	return e
}

func (e *ExecutorMock) SetOutputPrefix(prefix string) *ExecutorMock {
	e.outputPrefix = prefix
	return e
}

func (e *ExecutorMock) AddOutputFilter(substring string) *ExecutorMock {
	e.filterSubstrings = append(e.filterSubstrings, substring)
	return e
}

func (e *ExecutorMock) Run() *Result {
	return buildResult("", nil)
}

func (e *ExecutorMock) Capture() *Result {
	return buildResult("", nil)
}

func (e *ExecutorMock) CaptureAndTrim() *Result {
	return buildResult("", nil)
}
