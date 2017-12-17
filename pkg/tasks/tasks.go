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
	Load(map[interface{}]interface{}) (bool, error)
	Perform() error
}

func BuildFromDefinition(definition map[interface{}]interface{}) (task Task, err error) {
	task = &Custom{}
	ok, err := task.Load(definition)
	if err != nil {
		return nil, err
	}
	if ok {
		return task, nil
	}

	task = &Pip{}
	ok, err = task.Load(definition)
	if err != nil {
		return nil, err
	}
	if ok {
		return task, nil
	}

	return nil, fmt.Errorf("unknown task: %+v", definition)
}
