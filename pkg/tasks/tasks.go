package tasks

import (
	"github.com/pior/dad/pkg/termui"
)

type Task interface {
	Load(interface{}) (bool, error)
	Perform(*termui.UI) error
}

type TaskWithFeature interface {
	Task
	Features() map[string]string
}
