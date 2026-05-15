package integration

import (
	"testing"

	testcontext "github.com/devbuddy/devbuddy/tests/context"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Worktree_New_And_Cd(t *testing.T) {
	c := CreateContextAndInit(t)
	projectPath := "/home/tester/src/github.com/orgname/projname"
	worktreePath := "/home/tester/src/github.com/orgname/projname--feature-a"

	createGitProject(t, c, projectPath)
	c.Cd(t, projectPath)

	output := c.Run(t, "bud wt new feature-a")
	OutputContains(t, output, "created worktree", "feature-a", worktreePath)
	c.AssertExist(t, worktreePath+"/dev.yml")

	branch := c.Run(t, "git -C "+worktreePath+" branch --show-current")
	require.Equal(t, []string{"feature-a"}, branch)

	output = c.Run(t, "bud cd feature-a")
	OutputContains(t, output, "jumping to", "projname--feature-a")
	require.Equal(t, worktreePath, c.Cwd(t))

	c.Cd(t, projectPath)
	output = c.Run(t, "bud wt cd feature-a")
	OutputContains(t, output, "jumping to", "feature-a")
	require.Equal(t, worktreePath, c.Cwd(t))
}

func Test_Cmd_Worktree_BranchConflict(t *testing.T) {
	c := CreateContextAndInit(t)
	projectPath := "/home/tester/src/github.com/orgname/projname"
	worktreePath := "/home/tester/src/github.com/orgname/projname--feature-a"

	createGitProject(t, c, projectPath)
	c.Cd(t, projectPath)
	c.Run(t, "bud wt new feature-a")

	output := c.Run(t, "bud wt new duplicate feature-a", testcontext.ExitCode(1))
	OutputContains(t, output, "branch feature-a is already checked out", worktreePath, "bud wt cd feature-a")
}

func createGitProject(t *testing.T, c *testcontext.TestContext, path string) {
	t.Helper()

	c.Run(t, "mkdir -p "+path)
	c.Run(t, "git -C "+path+" init")
	c.Run(t, "git -C "+path+" checkout -b main")
	c.Run(t, "git -C "+path+" config user.email tester@example.com")
	c.Run(t, "git -C "+path+" config user.name Tester")
	c.WriteLines(t, path+"/dev.yml", "commands: {}")
	c.Run(t, "git -C "+path+" add dev.yml")
	c.Run(t, "git -C "+path+" commit -m init")
}
