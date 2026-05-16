package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Init(t *testing.T) {
	c := harness.NewDockerPTYInit(t)

	output := c.Run(t, "bud init default")
	require.Equal(t, []string{
		"🐼  Created dev.yml with template default",
		"⚠️   Open dev.yml to adjust for your needs.",
		"🐼  activated: env",
	}, output)

	devFile := c.Cat(t, "dev.yml")
	require.Contains(t, devFile, "# DevBuddy config file")

	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(devFile), &data)
	require.NoError(t, err)

	require.Contains(t, data, "up")
	require.Contains(t, data, "commands")
	require.Contains(t, data, "env")
	require.Contains(t, data, "open")
}
