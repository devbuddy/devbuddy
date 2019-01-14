package executor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewExecutor(t *testing.T) {
	result := NewBuilder()

	executor := result.NewExecutor("false")
	require.Equal(t, executor, New("false"))
}

func TestNewShell(t *testing.T) {
	result := NewBuilder()

	executor := result.NewShell("true")
	require.Equal(t, executor, NewShell("true"))
}
