package tasks

import (
	"github.com/pior/dad/pkg/project"
)

type Task interface {
	Load(interface{}) (bool, error)
	Perform(*Context) error
}

type TaskWithFeature interface {
	Task
	Feature(*project.Project) (string, string)
}
