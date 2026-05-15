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

	output := c.Run(t, "bud tree new feature-a")
	OutputContains(t, output, "created worktree", "feature-a", worktreePath)
	c.AssertExist(t, worktreePath+"/dev.yml")

	branch := c.Run(t, "git -C "+worktreePath+" branch --show-current")
	require.Equal(t, []string{"feature-a"}, branch)

	c.Write(t, worktreePath+"/scratch.txt", "dirty\n")
	output = c.Run(t, "bud tree list feature-a")
	OutputContains(t, output, "feature-a", "dirty", worktreePath)

	output = c.Run(t, "bud cd feature-a")
	OutputContains(t, output, "jumping to", "projname--feature-a")
	require.Equal(t, worktreePath, c.Cwd(t))

	c.Cd(t, projectPath)
	output = c.Run(t, "bud tree cd feature-a")
	OutputContains(t, output, "jumping to", "feature-a")
	require.Equal(t, worktreePath, c.Cwd(t))

	c.Run(t, "bud where")
	c.AssertContains(t, worktreePath+"/pwd.txt", worktreePath)
}

func Test_Cmd_Tree_Cd_Matches_Branch_When_Different_From_Worktree_Name(t *testing.T) {
	c := CreateContextAndInit(t)
	projectPath := "/home/tester/src/github.com/orgname/projname"
	worktreePath := "/home/tester/src/github.com/orgname/projname--agent-1"

	createGitProject(t, c, projectPath)
	c.Cd(t, projectPath)
	c.Run(t, "bud tree new agent-1 feature-a")

	output := c.Run(t, "bud cd feature-a")
	OutputContains(t, output, "jumping to", "feature-a")
	require.Equal(t, worktreePath, c.Cwd(t))
}

func Test_Cmd_Tree_Switch_Interactive(t *testing.T) {
	c := CreateContextAndInit(t)
	projectPath := "/home/tester/src/github.com/orgname/projname"
	worktreePath := "/home/tester/src/github.com/orgname/projname--feature-a"

	createGitProject(t, c, projectPath)
	c.Cd(t, projectPath)
	c.Run(t, "bud tree new feature-a")

	c.Send(t, "bud tree switch\n")
	c.Expect(t, "Select worktree")
	c.Send(t, "\x1b[B\r")
	c.WaitPrompt(t)

	require.Equal(t, worktreePath, c.Cwd(t))
}

func Test_Cmd_Tree_BranchConflict(t *testing.T) {
	c := CreateContextAndInit(t)
	projectPath := "/home/tester/src/github.com/orgname/projname"
	worktreePath := "/home/tester/src/github.com/orgname/projname--feature-a"

	createGitProject(t, c, projectPath)
	c.Cd(t, projectPath)
	c.Run(t, "bud tree new feature-a")

	output := c.Run(t, "bud tree new duplicate feature-a", testcontext.ExitCode(1))
	OutputContains(t, output, "branch feature-a is already checked out", worktreePath, "bud tree cd feature-a")
}

func createGitProject(t *testing.T, c *testcontext.TestContext, path string) {
	t.Helper()

	c.Run(t, "mkdir -p "+path)
	c.Run(t, "git -C "+path+" init")
	c.Run(t, "git -C "+path+" checkout -b main")
	c.Run(t, "git -C "+path+" config user.email tester@example.com")
	c.Run(t, "git -C "+path+" config user.name Tester")
	c.WriteLines(t, path+"/dev.yml",
		"commands:",
		"  where: pwd > pwd.txt",
	)
	c.Run(t, "git -C "+path+" add dev.yml")
	c.Run(t, "git -C "+path+" commit -m init")
}
