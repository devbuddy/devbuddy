package executor

type executorBuilderImpl struct {
}

// NewBuilder returns an *ExecutorBuilderImpl that will build Executor
func NewBuilder() ExecutorBuilder {
	return &executorBuilderImpl{}
}

// NewExecutor returns an Executor that will run the program with arguments
func (e *executorBuilderImpl) NewExecutor(program string, args ...string) Executor {
	return New(program, args...)
}

// NewShell returns an Executor that will run the command line in a shell
func (e *executorBuilderImpl) NewShell(cmdline string) Executor {
	return NewShell(cmdline)
}
