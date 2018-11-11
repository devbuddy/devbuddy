package tasks

import (
	"fmt"
	"reflect"
)

type taskConfig struct {
	name    string
	payload interface{}
}

func (c *taskConfig) getStringProperty(name string, allowSingle bool) (string, error) {
	return c.getStringPropertyDefault(name, "", allowSingle)
}

func (c *taskConfig) getStringPropertyDefault(name string, defaultValue string, allowSingle bool) (string, error) {
	if allowSingle {
		if value, ok := c.payload.(string); ok {
			return value, nil
		}
	}

	properties, ok := c.payload.(map[interface{}]interface{})
	if !ok {
		message := "not a hash"
		if allowSingle {
			message = "not a string"
		}
		return "", fmt.Errorf("%s: %T (%+v)", message, c.payload, c.payload)
	}

	value, ok := properties[name]
	if !ok {
		if defaultValue != "" {
			return defaultValue, nil
		}
		return "", fmt.Errorf("missing key '%s'", name)
	}

	str, err := asString(value)
	if err != nil {
		return "", fmt.Errorf("%s: %T (%+v)", err, value, value)
	}

	return str, nil
}

func (c *taskConfig) getListOfStrings() ([]string, error) {
	strings, ok := c.payload.([]string)
	if ok {
		return strings, nil
	}

	return nil, fmt.Errorf("not a list of strings: type %T (\"%+v\")", c.payload, c.payload)
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
