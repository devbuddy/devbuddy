package tasks

import (
	"fmt"
	"reflect"
)

type taskConfig struct {
	name    string
	payload interface{}
}

func (c *taskConfig) getPayloadAsString() (string, error) {
	value, ok := c.payload.(string)
	if !ok {
		return "", fmt.Errorf("need a string, found: %T (%v)", c.payload, c.payload)
	}
	return value, nil
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

	return nil, fmt.Errorf("invalid task: \"%+v\"", definition)
}
