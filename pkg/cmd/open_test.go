package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenCmdUsageAndArgs(t *testing.T) {
	require.Equal(t, "open [pattern]", openCmd.Use)

	require.NoError(t, openCmd.Args(openCmd, []string{}))
	require.NoError(t, openCmd.Args(openCmd, []string{"docs"}))
	require.Error(t, openCmd.Args(openCmd, []string{"a", "b"}))
}
