package worktree

import (
	"github.com/devbuddy/devbuddy/pkg/executor"
)

func List(exec *executor.Executor, repoPath string) ([]Worktree, error) {
	cmd := executor.New("git", "worktree", "list", "--porcelain")
	cmd.Cwd = repoPath

	result := exec.Capture(cmd)
	if result.Error != nil {
		return nil, result.Error
	}

	return ParseListPorcelain(result.Output)
}

func AddNewBranch(exec *executor.Executor, repoPath, path, branch string) error {
	cmd := executor.New("git", "worktree", "add", "-b", branch, path)
	cmd.Cwd = repoPath
	return exec.Run(cmd).Error
}

func AddExistingBranch(exec *executor.Executor, repoPath, path, branch string) error {
	cmd := executor.New("git", "worktree", "add", path, branch)
	cmd.Cwd = repoPath
	return exec.Run(cmd).Error
}

func Remove(exec *executor.Executor, repoPath, path string) error {
	cmd := executor.New("git", "worktree", "remove", path)
	cmd.Cwd = repoPath
	return exec.Run(cmd).Error
}

func Prune(exec *executor.Executor, repoPath string) error {
	cmd := executor.New("git", "worktree", "prune")
	cmd.Cwd = repoPath
	return exec.Run(cmd).Error
}

func BranchExists(exec *executor.Executor, repoPath, branch string) (bool, error) {
	cmd := executor.New("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	cmd.Cwd = repoPath

	result := exec.Run(cmd)
	if result.Code == 1 {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
