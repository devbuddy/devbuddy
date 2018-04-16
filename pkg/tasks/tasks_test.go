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
