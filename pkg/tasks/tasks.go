package tasks

import (
	"errors"
	"fmt"
)

var taskFailed error

func init() {
	taskFailed = errors.New("task failed")
}

type Task interface {
	Load(interface{}) (bool, error)
	Perform() error
}

func BuildFromDefinition(definition interface{}) (task Task, err error) {
	task = &Custom{}
	ok, err := task.Load(definition)
	if err != nil {
		return nil, err
	}
	if ok {
		return
	}

	task = &Pip{}
	ok, err = task.Load(definition)
	if err != nil {
		return nil, err
	}
	if ok {
		return
	}

	task = &Python{}
	ok, err = task.Load(definition)
	if err != nil {
		return nil, err
	}
	if ok {
		return
	}

	task = &Unknown{}
	ok, err = task.Load(definition)
	if err != nil {
		return nil, err
	}
	if ok {
		return
	}

	return nil, fmt.Errorf("error parsing tasks: %+v", definition)
}
