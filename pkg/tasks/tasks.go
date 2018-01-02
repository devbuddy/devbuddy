package tasks

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/pior/dad/pkg/termui"
)

var taskFailed error
var allTasks = make(map[string]TaskBuilder)

func init() {
	taskFailed = errors.New("task failed")
}

type Task interface {
	Load(interface{}) (bool, error)
	Perform(*termui.UI) error
}

type TaskWithFeature interface {
	Task
	Features() map[string]string
}

type TaskBuilder func() Task

func BuildFromDefinition(definition interface{}) (task Task, err error) {
	taskBuilder, err := findTaskBuilder(definition)
	if err != nil {
		return nil, err
	}
	task = taskBuilder()
	ok, err := task.Load(definition)
	if err != nil {
		return nil, err
	}
	if ok {
		return task, nil
	}

	return nil, fmt.Errorf("error parsing tasks: %+v", definition)
}

func findTaskBuilder(definition interface{}) (TaskBuilder, error) {
	name, err := extractTaskName(definition)
	if err != nil {
		return nil, fmt.Errorf("%s (%+v)", err, definition)
	}
	taskBuilder, found := allTasks[name]
	if found {
		return taskBuilder, nil
	}
	return NewUnknown, nil
}

func extractTaskName(definition interface{}) (string, error) {
	val := reflect.ValueOf(definition)
	if val.Kind() == reflect.Map {
		keys := val.MapKeys()
		if len(keys) != 1 {
			return "", fmt.Errorf("invalid map length")
		}
		definition = keys[0].Interface()
	}

	if name, ok := definition.(string); ok {
		return name, nil
	}

	return "", fmt.Errorf("invalid structure")
}
