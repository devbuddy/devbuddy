package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("env", "Env", func(_ *api.TaskConfig, _ *api.Task) error {
		panic("env tasks are constructed directly via newEnvTask in task_list.go")
	})
}
