package executor

// Executor prepares and run a command execution
type Executor interface {
	SetCwd(cwd string) Executor
	SetEnv(env []string) Executor
	SetEnvVar(name, value string) Executor
	SetOutputPrefix(prefix string) Executor
	SetPTY(enabled bool) Executor
	AddOutputFilter(substring string) Executor
	Run() *Result
	Capture() *Result
	CaptureAndTrim() *Result
}
