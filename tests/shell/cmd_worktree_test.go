package integration

import (
	"strings"
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
	require.Equal(t, worktreePath, c.Cwd(t))

	branch := c.Run(t, "git -C "+worktreePath+" branch --show-current")
	require.Equal(t, []string{"feature-a"}, branch)
	c.Run(t, "bud tree new agent-1 long-feature-branch")

	c.Write(t, worktreePath+"/scratch.txt", "dirty\n")
	output = c.Run(t, "bud tree list")
	OutputContains(t, output, "BRANCH", "HEAD", "STATE", "MODIFIED", "PATH")
	OutputContains(t, output, "feature-a", "dirty", worktreePath)
	assertPathColumnAligned(t, output)

	output = c.Run(t, `bud __complete tree cd ""`)
	OutputContains(t, output, "feature-a", "long-feature-branch")

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
	c := CreatePTYContextAndInit(t)
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

func Test_Cmd_Tree_Prune_Asks_To_Delete_Inactive_Worktrees(t *testing.T) {
	c := CreatePTYContextAndInit(t)
	projectPath := "/home/tester/src/github.com/orgname/projname"
	oldWorktreePath := "/home/tester/src/github.com/orgname/projname--old-branch"
	recentWorktreePath := "/home/tester/src/github.com/orgname/projname--recent-branch"

	createGitProject(t, c, projectPath)
	c.Cd(t, projectPath)
	c.Run(t, "bud tree new old-branch")
	c.Cd(t, projectPath)
	c.Run(t, "bud tree new recent-branch")
	c.Cd(t, projectPath)
	c.Run(t, "touch -d '8 days ago' "+oldWorktreePath)

	c.Send(t, "bud tree prune\n")
	c.Expect(t, "Delete inactive worktree old-branch")
	c.Send(t, "y\r")
	c.WaitPrompt(t)

	c.Run(t, "test ! -e "+oldWorktreePath)
	c.AssertExist(t, recentWorktreePath)
}

func Test_Cmd_Tree_Has_No_Short_Aliases(t *testing.T) {
	c := CreateContext(t)

	output := c.Run(t, "bud wt", testcontext.ExitCode(1))
	OutputContains(t, output, `unknown command "wt"`)

	output = c.Run(t, "bud worktree", testcontext.ExitCode(1))
	OutputContains(t, output, `unknown command "worktree"`)
}

func assertPathColumnAligned(t *testing.T, lines []string) {
	t.Helper()

	pathColumn := -1
	for _, line := range lines {
		if !strings.Contains(line, "/home/tester/src") {
			continue
		}

		column := strings.Index(line, "/home/tester/src")
		if pathColumn == -1 {
			pathColumn = column
			continue
		}
		require.Equal(t, pathColumn, column, "path column should be aligned in line %q", line)
	}
	require.NotEqual(t, -1, pathColumn, "expected at least one worktree path")
}

func createGitProject(t *testing.T, c *testcontext.TestContext, path string) {
	t.Helper()

	// Write via host path first so the directory is created with host-side
	// ownership before the container touches it (avoids chmod permission errors).
	c.WriteLines(t, path+"/dev.yml",
		"commands:",
		"  where: pwd > pwd.txt",
	)
	c.Run(t, "git -C "+path+" init")
	c.Run(t, "git -C "+path+" checkout -b main")
	c.Run(t, "git -C "+path+" config user.email tester@example.com")
	c.Run(t, "git -C "+path+" config user.name Tester")
	c.Run(t, "git -C "+path+" add dev.yml")
	c.Run(t, "git -C "+path+" commit -m init")
}
