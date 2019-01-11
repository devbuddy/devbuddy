package taskapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryUnknown(t *testing.T) {
	task := &Task{}
	parseUnknown(&TaskConfig{"sometask", "somevalue"}, task)

	// require.Equal(t, "nopenope", task.name)
	require.Equal(t, "", task.header)
}
