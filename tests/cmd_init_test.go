package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func Test_Cmd_Init(t *testing.T) {
	c := CreateContextAndInit(t)

	output := c.Run(t, "bud init default")
	require.Equal(t, []string{
		"🐼  Created dev.yml with template default",
		"⚠️   Open dev.yml to adjust for your needs.",
		"🐼  env activated.",
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
