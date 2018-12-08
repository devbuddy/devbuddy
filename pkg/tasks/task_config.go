package tasks

import (
	"fmt"
	"reflect"
)

type taskConfig struct {
	name    string
	payload interface{}
}

type propertyNotFoundError struct {
	name string
}

func (e propertyNotFoundError) Error() string {
	return fmt.Sprintf("property \"%s\" not found", e.name)
}

func (c *taskConfig) getListOfStrings() ([]string, error) {
	return asListOfStrings(c.payload)
}

func (c *taskConfig) getProperty(name string) (interface{}, error) {
	if c.payload == nil {
		return nil, propertyNotFoundError{name: name}
	}

	properties, ok := c.payload.(map[interface{}]interface{})
	if !ok {
		return "", fmt.Errorf("not a hash: %T (%+v)", c.payload, c.payload)
	}

	value, present := properties[name]
	if present {
		return value, nil
	}
	return nil, propertyNotFoundError{name: name}
}

func (c *taskConfig) getPropertyDefault(name string, defaultValue interface{}) (interface{}, error) {
	value, err := c.getProperty(name)
	if err == nil {
		return value, nil
	}
	if _, ok := err.(propertyNotFoundError); ok {
		return defaultValue, nil
	}
	return nil, err
}

func (c *taskConfig) getStringPropertyAllowSingle(name string) (string, error) {
	if value, ok := c.payload.(string); ok {
		return value, nil
	}
	return c.getStringProperty(name)
}

func (c *taskConfig) getStringProperty(name string) (string, error) {
	value, err := c.getProperty(name)
	if err != nil {
		return "", err
	}
	return asString(value)
}

func (c *taskConfig) getStringPropertyDefault(name string, defaultValue string) (string, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return "", err
	}
	return asString(value)
}

func (c *taskConfig) getListOfStringsPropertyDefault(name string, defaultValue []string) ([]string, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return nil, err
	}
	return asListOfStrings(value)
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
