package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryUnknown(t *testing.T) {
	task, err := loadTestTask(t, `nopenope`)

	require.NoError(t, err)
	require.NotNil(t, task)

	require.Equal(t, "nopenope", task.name())
	require.Equal(t, "", task.header())
}
