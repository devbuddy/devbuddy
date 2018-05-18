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

	require.Equal(t, "custom-command", task.(*Custom).command)
	require.Equal(t, "test-command", task.(*Custom).condition)
}
func TestCustomName(t *testing.T) {
	task := ensureLoadTestTask(t, `
custom:
  name: NAMENAME
  met?: test-command
  meet: custom-command
`)

	require.Equal(t, "NAMENAME", task.header())
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
