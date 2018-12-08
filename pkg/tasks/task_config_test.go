package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskConfigProperty(t *testing.T) {
	config := &taskConfig{
		name:    "test",
		payload: map[interface{}]interface{}{"key": []interface{}{"one", "two"}},
	}

	val, err := config.getListOfStringsPropertyDefault("key", []string{})
	require.NoError(t, err)
	require.Equal(t, []string{"one", "two"}, val)

	val, err = config.getListOfStringsPropertyDefault("nope", []string{"three"})
	require.NoError(t, err)
	require.Equal(t, []string{"three"}, val)
}

func TestTaskConfigPropertyInvalid(t *testing.T) {
	config := &taskConfig{name: "test", payload: "a string"}

	_, err := config.getListOfStringsPropertyDefault("key", []string{})
	require.Error(t, err)
	require.Equal(t, "not a hash: string (a string)", err.Error())
}

func TestTaskConfigStringOrStringProperty(t *testing.T) {
	config := &taskConfig{name: "test", payload: "value"}

	val, err := config.getStringPropertyAllowSingle("key")
	require.NoError(t, err)
	require.Equal(t, val, "value")
}

func TestTaskConfigStringProperty(t *testing.T) {
	value := map[interface{}]interface{}{"key": "val"}
	config := &taskConfig{name: "test", payload: value}

	val, err := config.getStringProperty("key")
	require.NoError(t, err)
	require.Equal(t, val, "val")

	val, err = config.getStringProperty("nope")
	require.Error(t, err)
	require.Equal(t, "property \"nope\" not found", err.Error())
}

func TestTaskConfigStringPropertyInvalid(t *testing.T) {
	config := &taskConfig{name: "test", payload: 42}
	_, err := config.getStringPropertyAllowSingle("key")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a hash")
	require.Contains(t, err.Error(), "int (42)")

	config = &taskConfig{name: "test", payload: false}
	_, err = config.getStringPropertyAllowSingle("key")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a hash")
	require.Contains(t, err.Error(), "bool (false)")

	config = &taskConfig{name: "test", payload: "thisisastring"}
	_, err = config.getStringProperty("key1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a hash")
	require.Contains(t, err.Error(), "thisisastring")

	payload := map[interface{}]interface{}{"version": 3.6}
	config = &taskConfig{name: "test", payload: payload}
	_, err = config.getStringProperty("version")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a string")
	require.Contains(t, err.Error(), "float64 (3.6)")
}

func TestTaskConfigListOfStrings(t *testing.T) {
	value := []interface{}{"one", "two"}
	config := &taskConfig{name: "test", payload: value}

	result, err := config.getListOfStrings()
	require.NoError(t, err)
	require.Equal(t, []string{"one", "two"}, result)
}

func TestTaskConfigListOfStringsEmpty(t *testing.T) {
	config := &taskConfig{name: "test", payload: []interface{}{}}

	result, err := config.getListOfStrings()
	require.NoError(t, err)
	require.Equal(t, []string{}, result)
}

func TestTaskConfigListOfStringsInvalidElement(t *testing.T) {
	config := &taskConfig{name: "test", payload: []interface{}{"one", 2}}

	_, err := config.getListOfStrings()
	require.Error(t, err)
	require.Equal(t, "not a list of strings: invalid element: type int (2)", err.Error())
}

func TestTaskConfigListOfStringsInvalidType(t *testing.T) {
	config := &taskConfig{name: "test", payload: "plop"}

	_, err := config.getListOfStrings()
	require.Error(t, err)
	require.Equal(t, "not a list of strings: type string (plop)", err.Error())
}
