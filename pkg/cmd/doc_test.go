package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootDocFlagPrintsAgentGuide(t *testing.T) {
	rootCmd := build("test-version")
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetArgs([]string{"--doc"})

	err := rootCmd.Execute()

	require.NoError(t, err)
	require.Contains(t, output.String(), "# DevBuddy Project Guide")
	require.Contains(t, output.String(), "dev.yml")
	require.Contains(t, output.String(), "bud up")
	require.Contains(t, output.String(), "bud --help")
	require.Contains(t, output.String(), "bud <command>")
	require.Contains(t, output.String(), "bud --shell-init")
}
