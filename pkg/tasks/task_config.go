package tasks

import (
	"fmt"
)

type propertyNotFoundError struct {
	name string
}

func (e propertyNotFoundError) Error() string {
	return fmt.Sprintf("property \"%s\" not found", e.name)
}

// TaskConfig represents a task as defined in dev.yml
type TaskConfig struct {
	name    string
	payload interface{}
}

// GetListOfStrings expects the payload to be a list of string, returns it.
//
// YAML Example:
// - taskname:
//     - foo
//     - bar
func (c *TaskConfig) GetListOfStrings() ([]string, error) {
	return asListOfStrings(c.payload)
}

// GetListOfStringsPropertyDefault expects the payload to be a list of string, returns it.
// Or returns the default value specified if the payload is nil.
//
// YAML Example:
// - taskname:
//     - foo
//     - bar
func (c *TaskConfig) GetListOfStringsPropertyDefault(name string, defaultValue []string) ([]string, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return nil, err
	}
	return asListOfStrings(value)
}

// GetStringPropertyAllowSingle expects one of those two situations:
// 1. the payload is a string, returns it directly
// 2. the payload is a hash of string, returns the value for the name specified
//
// YAML Example:
// - taskname: foo
// or
// - taskname:
//     <name>: foo
func (c *TaskConfig) GetStringPropertyAllowSingle(name string) (string, error) {
	if value, ok := c.payload.(string); ok {
		return value, nil
	}
	return c.GetStringProperty(name)
}

// GetStringProperty expects the payload to be a hash of string, returns the value for the name specified.
//
// YAML Example:
// - taskname:
//     <name>: foo
func (c *TaskConfig) GetStringProperty(name string) (string, error) {
	value, err := c.getProperty(name)
	if err != nil {
		return "", err
	}
	return asString(value)
}

// GetStringPropertyDefault expects the payload to be a hash of string, returns the value for the name specified.
// Or returns the default value specified if the payload is nil.
//
// YAML Example:
// - taskname:
//     <name>: foo
func (c *TaskConfig) GetStringPropertyDefault(name string, defaultValue string) (string, error) {
	value, err := c.getPropertyDefault(name, defaultValue)
	if err != nil {
		return "", err
	}
	return asString(value)
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
