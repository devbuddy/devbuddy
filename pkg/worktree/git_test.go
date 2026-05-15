package worktree

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/executor"
)

type runnerSpy struct {
	runCmds     []*executor.Command
	captureCmds []*executor.Command
	output      string
}

func (s *runnerSpy) Run(cmd *executor.Command) *executor.Result {
	s.runCmds = append(s.runCmds, cmd)
	return &executor.Result{}
}

func (s *runnerSpy) Capture(cmd *executor.Command) *executor.Result {
	s.captureCmds = append(s.captureCmds, cmd)
	return &executor.Result{Output: s.output}
}

func TestListRunsGitWorktreeListPorcelain(t *testing.T) {
	runner := &runnerSpy{
		output: "worktree /src/github.com/acme/api\nHEAD 1111111111111111111111111111111111111111\nbranch refs/heads/main\n",
	}
	exec := &executor.Executor{Runner: runner}

	got, err := List(exec, "/src/github.com/acme/api")

	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "main", got[0].Branch)
	require.Len(t, runner.captureCmds, 1)
	require.Equal(t, "git", runner.captureCmds[0].Program)
	require.Equal(t, []string{"worktree", "list", "--porcelain"}, runner.captureCmds[0].Args)
	require.Equal(t, "/src/github.com/acme/api", runner.captureCmds[0].Cwd)
}

func TestAddNewBranchBuildsGitCommand(t *testing.T) {
	runner := &runnerSpy{}
	exec := &executor.Executor{Runner: runner}

	err := AddNewBranch(exec, "/src/github.com/acme/api", "/src/github.com/acme/api--feature-a", "feature-a")

	require.NoError(t, err)
	require.Len(t, runner.runCmds, 1)
	require.Equal(t, "git", runner.runCmds[0].Program)
	require.Equal(t, []string{"worktree", "add", "-b", "feature-a", "/src/github.com/acme/api--feature-a"}, runner.runCmds[0].Args)
	require.Equal(t, "/src/github.com/acme/api", runner.runCmds[0].Cwd)
}

func TestAddExistingBranchBuildsGitCommand(t *testing.T) {
	runner := &runnerSpy{}
	exec := &executor.Executor{Runner: runner}

	err := AddExistingBranch(exec, "/src/github.com/acme/api", "/src/github.com/acme/api--feature-a", "feature-a")

	require.NoError(t, err)
	require.Len(t, runner.runCmds, 1)
	require.Equal(t, []string{"worktree", "add", "/src/github.com/acme/api--feature-a", "feature-a"}, runner.runCmds[0].Args)
	require.Equal(t, "/src/github.com/acme/api", runner.runCmds[0].Cwd)
}

func TestIsDirtyRunsGitStatusShort(t *testing.T) {
	runner := &runnerSpy{output: " M README.md\n"}
	exec := &executor.Executor{Runner: runner}

	got, err := IsDirty(exec, "/src/github.com/acme/api--feature-a")

	require.NoError(t, err)
	require.True(t, got)
	require.Len(t, runner.captureCmds, 1)
	require.Equal(t, "git", runner.captureCmds[0].Program)
	require.Equal(t, []string{"status", "--short"}, runner.captureCmds[0].Args)
	require.Equal(t, "/src/github.com/acme/api--feature-a", runner.captureCmds[0].Cwd)
}
