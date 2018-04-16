package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryInvalid(t *testing.T) {
	task, err := loadTestTask(t, `42`)

	require.NoError(t, err)
	require.NotNil(t, task)

	require.Equal(t, "Invalid task", task.name())
	require.Equal(t, "", task.header())
}
