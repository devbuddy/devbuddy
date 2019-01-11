package taskapi

import (
	"fmt"
	"reflect"
)

type TaskConfig struct {
	name    string
	payload interface{}
}

type propertyNotFoundError struct {
	name string
}

func (e propertyNotFoundError) Error() string {
	return fmt.Sprintf("property \"%s\" not found", e.name)
}

func (c *TaskConfig) GetListOfStrings() ([]string, error) {
	return asListOfStrings(c.payload)
}

func (c *TaskConfig) getProperty(name string) (interface{}, error) {
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

func (c *TaskConfig) getPropertyDefault(name string, defaultValue interface{}) (interface{}, error) {
	value, err := c.getProperty(name)
	if err == nil {
		return value, nil
	}
	if _, ok := err.(propertyNotFoundError); ok {
		return defaultValue, nil
	}
	return nil, err
}

func (c *TaskConfig) GetStringPropertyAllowSingle(name string) (string, error) {
	if value, ok := c.payload.(string); ok {
		return value, nil
	}
	return c.GetStringProperty(name)
}

func (c *TaskConfig) GetStringProperty(name string) (string, error) {
	value, err := c.getProperty(name)
	if err != nil {
		return "", err
	}
	return asString(value)
}

func (c *TaskConfig) GetStringPropertyDefault(name string, defaultValue string) (string, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return "", err
	}
	return asString(value)
}

func (c *TaskConfig) GetListOfStringsPropertyDefault(name string, defaultValue []string) ([]string, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return nil, err
	}
	return asListOfStrings(value)
}

func parseTaskConfig(definition interface{}) (*TaskConfig, error) {
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
		return &TaskConfig{name: name, payload: payload}, nil
	}

	if val.Kind() == reflect.String {
		return &TaskConfig{name: definition.(string), payload: nil}, nil
	}

	return nil, fmt.Errorf("invalid task: \"%+v\"", definition)
}
