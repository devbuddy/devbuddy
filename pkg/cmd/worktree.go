package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/integration"
	"github.com/devbuddy/devbuddy/pkg/worktree"
)

var worktreeCmd = &cobra.Command{
	Use:          "wt",
	Aliases:      []string{"worktree"},
	Short:        "Manage git worktrees",
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

var worktreeListCmd = &cobra.Command{
	Use:          "list [QUERY]",
	Short:        "List git worktrees",
	Args:         zeroOrOneArg,
	RunE:         worktreeListRun,
	SilenceUsage: true,
}

var worktreeNewCmd = &cobra.Command{
	Use:          "new NAME [BRANCH]",
	Short:        "Create a sibling git worktree",
	Args:         oneOrTwoArgs,
	RunE:         worktreeNewRun,
	SilenceUsage: true,
}

var worktreeCdCmd = &cobra.Command{
	Use:          "cd QUERY",
	Short:        "Jump to a git worktree",
	Args:         onlyOneArg,
	RunE:         worktreeCdRun,
	SilenceUsage: true,
}

var worktreeRemoveCmd = &cobra.Command{
	Use:          "remove QUERY",
	Short:        "Remove a git worktree",
	Args:         onlyOneArg,
	RunE:         worktreeRemoveRun,
	SilenceUsage: true,
}

var worktreePruneCmd = &cobra.Command{
	Use:          "prune",
	Short:        "Prune stale git worktree metadata",
	Args:         noArgs,
	RunE:         worktreePruneRun,
	SilenceUsage: true,
}

func init() {
	worktreeCmd.AddCommand(worktreeListCmd)
	worktreeCmd.AddCommand(worktreeNewCmd)
	worktreeCmd.AddCommand(worktreeCdCmd)
	worktreeCmd.AddCommand(worktreeRemoveCmd)
	worktreeCmd.AddCommand(worktreePruneCmd)
}

func worktreeListRun(_ *cobra.Command, args []string) error {
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	worktrees, err := worktree.List(ctx.Executor, ctx.Project.Path)
	if err != nil {
		return err
	}

	query := ""
	if len(args) == 1 {
		query = args[0]
	}

	for _, wt := range matchWorktrees(worktrees, query) {
		printWorktree(wt)
	}
	return nil
}

func worktreeNewRun(_ *cobra.Command, args []string) error {
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	worktrees, err := worktree.List(ctx.Executor, ctx.Project.Path)
	if err != nil {
		return err
	}

	name := args[0]
	branch := worktree.Slug(name)
	if len(args) == 2 {
		branch = args[1]
	}
	if branch == "" {
		return fmt.Errorf("worktree branch must contain letters or numbers")
	}

	if conflict := worktree.CheckedOutBranch(worktrees, branch); conflict != nil {
		return fmt.Errorf("branch %s is already checked out at %s\nRun: bud wt cd %s", branch, conflict.Path, branch)
	}

	repoPath := mainWorktreePath(worktrees, ctx.Project.Path)
	path, err := worktree.ManagedPath(repoPath, name)
	if err != nil {
		return err
	}

	branchExists, err := worktree.BranchExists(ctx.Executor, repoPath, branch)
	if err != nil {
		return err
	}

	if branchExists {
		err = worktree.AddExistingBranch(ctx.Executor, repoPath, path, branch)
	} else {
		err = worktree.AddNewBranch(ctx.Executor, repoPath, path, branch)
	}
	if err != nil {
		return err
	}

	fmt.Printf("🐼  created worktree %s at %s\n", branch, path)
	return nil
}

func worktreeCdRun(_ *cobra.Command, args []string) error {
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	wt, err := findWorktree(ctx.Executor, ctx.Project.Path, args[0])
	if err != nil {
		return err
	}

	ctx.UI.JumpProject(worktreeLabel(wt))
	return integration.AddFinalizerCd(wt.Path)
}

func worktreeRemoveRun(_ *cobra.Command, args []string) error {
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	worktrees, err := worktree.List(ctx.Executor, ctx.Project.Path)
	if err != nil {
		return err
	}
	matches := matchWorktrees(worktrees, args[0])
	if len(matches) == 0 {
		return fmt.Errorf("no worktree found for %s", args[0])
	}

	wt := matches[0]
	repoPath := mainWorktreePath(worktrees, ctx.Project.Path)
	if filepath.Clean(wt.Path) == filepath.Clean(repoPath) {
		return fmt.Errorf("refusing to remove the main worktree: %s", wt.Path)
	}

	if err := worktree.Remove(ctx.Executor, ctx.Project.Path, wt.Path); err != nil {
		return err
	}

	fmt.Printf("🐼  removed worktree %s\n", wt.Path)
	return nil
}

func worktreePruneRun(_ *cobra.Command, _ []string) error {
	ctx, err := context.LoadWithProject()
	if err != nil {
		return err
	}

	if err := worktree.Prune(ctx.Executor, ctx.Project.Path); err != nil {
		return err
	}

	fmt.Println("🐼  pruned stale worktree metadata")
	return nil
}

func oneOrTwoArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("expecting one or two arguments")
	}
	return nil
}

func findWorktree(exec *executor.Executor, repoPath, query string) (worktree.Worktree, error) {
	worktrees, err := worktree.List(exec, repoPath)
	if err != nil {
		return worktree.Worktree{}, err
	}

	matches := matchWorktrees(worktrees, query)
	if len(matches) == 0 {
		return worktree.Worktree{}, fmt.Errorf("no worktree found for %s", query)
	}
	return matches[0], nil
}

func matchWorktrees(worktrees []worktree.Worktree, query string) []worktree.Worktree {
	if query == "" {
		return worktrees
	}

	for _, wt := range worktrees {
		if worktreeExactMatch(wt, query) {
			return []worktree.Worktree{wt}
		}
	}

	index := make([]string, 0, len(worktrees))
	for _, wt := range worktrees {
		index = append(index, worktreeLabel(wt)+" "+filepath.Base(wt.Path))
	}

	matches := fuzzy.Find(query, index)
	found := make([]worktree.Worktree, 0, matches.Len())
	for _, match := range matches {
		found = append(found, worktrees[match.Index])
	}
	return found
}

func worktreeExactMatch(wt worktree.Worktree, query string) bool {
	base := filepath.Base(wt.Path)
	return wt.Branch == query || base == query || strings.HasSuffix(base, "--"+query)
}

func mainWorktreePath(worktrees []worktree.Worktree, fallback string) string {
	for _, wt := range worktrees {
		if !wt.Bare {
			return wt.Path
		}
	}
	return fallback
}

func printWorktree(wt worktree.Worktree) {
	head := wt.Head
	if len(head) > 7 {
		head = head[:7]
	}

	fmt.Printf("%s  %s  %s\n", worktreeLabel(wt), head, wt.Path)
}

func worktreeLabel(wt worktree.Worktree) string {
	if wt.Branch != "" {
		return wt.Branch
	}
	if wt.Detached {
		return "detached"
	}
	return strings.TrimPrefix(filepath.Base(wt.Path), filepath.Base(filepath.Dir(wt.Path))+"--")
}
