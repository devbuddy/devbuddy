package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvfileDefault(t *testing.T) {
	task := ensureLoadTestTask(t, `envfile`)

	require.Equal(t, "Task EnvFile () feature=envfile:.env actions=1", task.Describe())
}

func TestEnvfileCustomPath(t *testing.T) {
	task := ensureLoadTestTask(t, `envfile: config/local-dev.env`)

	require.Equal(t, "Task EnvFile () feature=envfile:config/local-dev.env actions=1", task.Describe())
}
