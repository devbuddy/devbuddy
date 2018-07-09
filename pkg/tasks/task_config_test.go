package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskConfigStringOk(t *testing.T) {
	config := &taskConfig{name: "test", payload: "value"}

	val, err := config.getPayloadAsString()
	require.NoError(t, err)
	require.Equal(t, val, "value")
}

func TestTaskConfigStringWithNumber(t *testing.T) {
	config := &taskConfig{name: "test", payload: 42}

	_, err := config.getPayloadAsString()
	require.EqualError(t, err, "need a string, found: int (42)")
}

func TestTaskConfigStringWithBoolean(t *testing.T) {
	config := &taskConfig{name: "test", payload: false}

	_, err := config.getPayloadAsString()
	require.EqualError(t, err, "need a string, found: bool (false)")
}

func TestTaskConfigMapOk(t *testing.T) {
	value := map[interface{}]interface{}{"key1": "val1", "key2": "val2"}

	config := &taskConfig{name: "test", payload: value}

	val, err := config.getPayloadAsStringMap()
	require.NoError(t, err)
	require.Equal(t, val, map[string]string{"key1": "val1", "key2": "val2"})
}

func TestTaskConfigMapNotAMap(t *testing.T) {
	config := &taskConfig{name: "test", payload: "thisisastring"}
	_, err := config.getPayloadAsStringMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a hash")
	require.Contains(t, err.Error(), "thisisastring")
}

func TestTaskConfigMapWithInvalidValues(t *testing.T) {
	payload := map[interface{}]interface{}{"version": 3.6}

	config := &taskConfig{name: "test", payload: payload}
	_, err := config.getPayloadAsStringMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "3.6")
	require.Contains(t, err.Error(), "version")
	require.Contains(t, err.Error(), "not a string")
}
