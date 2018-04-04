package tasks

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustom(t *testing.T) {
	task := ensureLoadTestTask(t, `
custom:
  met?: test-command
  meet: custom-command
`)

	require.Equal(t, task.(*Custom).command, "custom-command")
	require.Equal(t, task.(*Custom).condition, "test-command")
}

func TestCustomWithBoolean(t *testing.T) {
	_, err := loadTestTask(t, `
custom:
  met?: false
  meet: custom-command
`)

	require.Error(t, err, "buildFromDefinition() should have failed")
	if !strings.Contains(err.Error(), "boolean") {
		t.Fatal("invalid err")
	}
}
