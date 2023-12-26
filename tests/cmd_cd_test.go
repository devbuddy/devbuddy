package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Cd(t *testing.T) {
	c := CreateContext(t)

	c.Run(t, "mkdir -p /home/tester/src/github.com/orgname/projname")

	output := c.Run(t, `eval "$(bud --shell-init)"`)
	require.Len(t, output, 0)

	output = c.Run(t, "bud cd projname")
	require.Len(t, output, 1)
	require.Equal(t, output[0], "🐼  jumping to github.com:orgname/projname")

	cwd := c.Cwd(t)
	require.Equal(t, "/home/tester/src/github.com/orgname/projname", cwd)
}

func Test_Cmd_Cd_Matching(t *testing.T) {
	c := CreateContext(t)

	project1 := "/home/tester/src/github.com/devbuddy_tests/project"
	project2 := "/home/tester/src/github.com/devbuddy_tests/project2"

	c.Run(t, "mkdir -p "+project1)
	c.Run(t, "mkdir -p "+project2)

	output := c.Run(t, `eval "$(bud --shell-init)"`)
	require.Len(t, output, 0)

	tests := map[string]string{
		"devbuddy_tests/project":  project1,
		"devbuddy_tests/project2": project2,

		"devbuddyproject":  project1,
		"devbuddyproject2": project2,

		"proj": project1,
		"pro2": project2,

		"dtp":  project1,
		"dtp2": project2,
	}

	for test, projectPath := range tests {
		t.Run(test, func(t *testing.T) {
			c.Run(t, "bud cd "+test)

			cwd := c.Cwd(t)
			require.Equal(t, projectPath, cwd)
		})
	}
}
