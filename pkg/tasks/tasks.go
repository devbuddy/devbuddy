package tasks

import (
	"fmt"
	"reflect"

	"github.com/pior/dad/pkg/project"
)

type Task interface {
	Load(*taskConfig) (bool, error)
	Perform(*Context) error
}

type TaskWithFeature interface {
	Task
	Feature(*project.Project) (string, string)
}

type taskConfig struct {
	name    string
	payload interface{}
}

func parseTaskConfig(definition interface{}) (*taskConfig, error) {
	val := reflect.ValueOf(definition)

	if val.Kind() == reflect.Map {
		keys := val.MapKeys()
		if len(keys) != 1 {
			return nil, fmt.Errorf("invalid map length")
		}
		name, ok := keys[0].Interface().(string)
		if !ok {
			return nil, fmt.Errorf("task name should be a string")
		}
		payload := val.MapIndex(keys[0]).Interface()
		return &taskConfig{name: name, payload: payload}, nil
	}

	if val.Kind() == reflect.String {
		return &taskConfig{name: definition.(string), payload: nil}, nil
	}

	return nil, fmt.Errorf("invalid structure")
}
