package tasks

import ()

type Task interface {
	Load(interface{}) (bool, error)
	Perform(*Context) error
}

type TaskWithFeature interface {
	Task
	Features() map[string]string
}
