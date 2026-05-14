package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode(t *testing.T) {
	task := ensureLoadTestTask(t, `node: 20.11.1`)

	require.Equal(t, "Task NodeJS (20.11.1) feature=node:20.11.1 actions=2", task.Describe())
}

func TestNodeMissingVersion(t *testing.T) {
	_, err := loadTestTask(t, `
node:
  npm: true
`)

	require.EqualError(t, err, `task "node": property "version" not found`)
}

func TestNodeInvalidVersionType(t *testing.T) {
	_, err := loadTestTask(t, `
node:
  version: 20
`)

	require.EqualError(t, err, `task "node": key "version": expecting a string, found a int (20)`)
}
