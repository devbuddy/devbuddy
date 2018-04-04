package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithString(t *testing.T) {
	s, err := asString("poipoi")

	require.NoError(t, err, "should not err for a string")
	require.Equal(t, s, "poipoi")
}

func TestWithBoolean(t *testing.T) {
	_, err := asString(false)

	require.Error(t, err, "should err for a boolean")
}
