package executor

// ExecutorBuilder build executor with specific program and args
type ExecutorBuilder interface {
	NewExecutor(program string, args ...string) Executor
	NewShell(cmdline string) Executor
}
