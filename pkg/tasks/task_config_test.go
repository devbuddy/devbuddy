package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskConfigStringOk(t *testing.T) {
	config := &taskConfig{name: "test", payload: "value"}

	val, err := config.getStringProperty("key", true)
	require.NoError(t, err)
	require.Equal(t, val, "value")
}

func TestTaskConfigStringWithNumber(t *testing.T) {
	config := &taskConfig{name: "test", payload: 42}

	_, err := config.getStringProperty("key", true)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a string")
	require.Contains(t, err.Error(), "int (42)")
}

func TestTaskConfigStringWithBoolean(t *testing.T) {
	config := &taskConfig{name: "test", payload: false}

	_, err := config.getStringProperty("key", true)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a string")
	require.Contains(t, err.Error(), "bool (false)")
}

func TestTaskConfigMapOk(t *testing.T) {
	value := map[interface{}]interface{}{"key1": "val1", "key2": "val2"}
	config := &taskConfig{name: "test", payload: value}

	val, err := config.getStringProperty("key1", false)
	require.NoError(t, err)
	require.Equal(t, val, "val1")
}

func TestTaskConfigMapNotAMap(t *testing.T) {
	config := &taskConfig{name: "test", payload: "thisisastring"}

	_, err := config.getStringProperty("key1", false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a hash")
	require.Contains(t, err.Error(), "thisisastring")
}

func TestTaskConfigMapWithInvalidValues(t *testing.T) {
	payload := map[interface{}]interface{}{"version": 3.6}
	config := &taskConfig{name: "test", payload: payload}

	_, err := config.getStringProperty("version", false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a string")
	require.Contains(t, err.Error(), "3.6")
	require.Contains(t, err.Error(), "float64")
}
