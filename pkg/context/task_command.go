package context

import "github.com/devbuddy/devbuddy/pkg/executor"

// RunTaskCommand records the command shown in task output and runs that same
// inspectable command request.
func (ctx *Context) RunTaskCommand(cmd *executor.Command) *executor.Result {
	ctx.UI.TaskCommand(cmd.Program, cmd.Args...)
	return ctx.Executor.Run(cmd)
}
