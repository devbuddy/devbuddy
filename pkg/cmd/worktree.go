package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/integration"
	"github.com/devbuddy/devbuddy/pkg/worktree"
)

var worktreeCmd = &cobra.Command{
	Use:          "tree",
	Short:        "Manage git worktrees",
	GroupID:      "devbuddy",
	RunE:         subcommandOnlyRun,
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

var worktreeSwitchCmd = &cobra.Command{
	Use:               "switch [QUERY]",
	Short:             "Select and jump to a git worktree",
	Args:              zeroOrOneArg,
	RunE:              worktreeSwitchRun,
	ValidArgsFunction: worktreeSwitchCompletions,
	SilenceUsage:      true,
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
	worktreeCmd.AddCommand(worktreeSwitchCmd)
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

	rows := buildWorktreeRows(ctx.Executor, matchWorktrees(worktrees, query))
	for _, line := range formatWorktreeRows(rows, true) {
		fmt.Println(line)
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
		return fmt.Errorf("branch %s is already checked out at %s\nRun: bud tree switch %s", branch, conflict.Path, branch)
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
	return integration.AddFinalizerCd(path)
}

func worktreeSwitchRun(_ *cobra.Command, args []string) error {
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

	matches := matchWorktrees(worktrees, query)
	if len(matches) == 0 {
		return fmt.Errorf("no worktree found for %s", query)
	}

	wt := matches[0]
	if query == "" && len(matches) > 1 {
		selected, err := selectWorktree(ctx.Executor, matches)
		if err != nil {
			return err
		}
		wt = selected
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

	worktrees, err := worktree.List(ctx.Executor, ctx.Project.Path)
	if err != nil {
		return err
	}

	for _, wt := range inactiveWorktrees(worktrees, time.Now(), 7*24*time.Hour) {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Delete inactive worktree %s at %s", worktreeLabel(wt), wt.Path),
			IsConfirm: true,
		}

		if _, err := prompt.Run(); err != nil {
			if errors.Is(err, promptui.ErrAbort) {
				continue
			}
			return err
		}

		if err := worktree.Remove(ctx.Executor, ctx.Project.Path, wt.Path); err != nil {
			return err
		}
		fmt.Printf("🐼  removed inactive worktree %s\n", wt.Path)
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

type switchItem struct {
	Label    string
	Worktree worktree.Worktree
}

func selectWorktree(exec *executor.Executor, worktrees []worktree.Worktree) (worktree.Worktree, error) {
	rows := buildWorktreeRows(exec, worktrees)
	labels := formatWorktreeRows(rows, false)
	items := make([]switchItem, 0, len(rows))
	for i, row := range rows {
		items = append(items, switchItem{
			Label:    labels[i],
			Worktree: row.Worktree,
		})
	}

	prompt := promptui.Select{
		Label:        "Select worktree",
		Items:        items,
		HideSelected: true,
		Templates:    worktreeSwitchTemplates(),
	}

	index, _, err := prompt.Run()
	if err != nil {
		return worktree.Worktree{}, err
	}
	return items[index].Worktree, nil
}

func worktreeSwitchTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🐼 {{ .Label | cyan }}",
		Inactive: "   {{ .Label }}",
		Selected: "🐼 {{ .Label }}",
	}
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

func inactiveWorktrees(worktrees []worktree.Worktree, now time.Time, maxAge time.Duration) []worktree.Worktree {
	if len(worktrees) == 0 {
		return nil
	}

	mainPath := mainWorktreePath(worktrees, "")
	cutoff := now.Add(-maxAge)
	var inactive []worktree.Worktree
	for _, wt := range worktrees {
		if filepath.Clean(wt.Path) == filepath.Clean(mainPath) {
			continue
		}

		info, err := os.Stat(wt.Path)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			inactive = append(inactive, wt)
		}
	}
	return inactive
}

type worktreeRow struct {
	Worktree worktree.Worktree
	Branch   string
	Head     string
	State    string
	Modified string
	Path     string
}

func buildWorktreeRows(exec *executor.Executor, worktrees []worktree.Worktree) []worktreeRow {
	rows := make([]worktreeRow, 0, len(worktrees))
	for _, wt := range worktrees {
		rows = append(rows, buildWorktreeRow(exec, wt))
	}
	return rows
}

func buildWorktreeRow(exec *executor.Executor, wt worktree.Worktree) worktreeRow {
	head := wt.Head
	if len(head) > 7 {
		head = head[:7]
	}

	state := "clean"
	dirty, err := worktree.IsDirty(exec, wt.Path)
	if err == nil && dirty {
		state = "dirty"
	}

	modified := "unknown"
	if info, err := os.Stat(wt.Path); err == nil {
		modified = info.ModTime().Format("2006-01-02")
	}

	return worktreeRow{
		Worktree: wt,
		Branch:   worktreeLabel(wt),
		Head:     head,
		State:    state,
		Modified: modified,
		Path:     wt.Path,
	}
}

func formatWorktreeRows(rows []worktreeRow, includeHeader bool) []string {
	allRows := rows
	if includeHeader {
		allRows = append([]worktreeRow{{
			Branch:   "BRANCH",
			Head:     "HEAD",
			State:    "STATE",
			Modified: "MODIFIED",
			Path:     "PATH",
		}}, rows...)
	}

	branchWidth := len("BRANCH")
	headWidth := len("HEAD")
	stateWidth := len("STATE")
	modifiedWidth := len("MODIFIED")
	for _, row := range rows {
		branchWidth = max(branchWidth, len(row.Branch))
		headWidth = max(headWidth, len(row.Head))
		stateWidth = max(stateWidth, len(row.State))
		modifiedWidth = max(modifiedWidth, len(row.Modified))
	}

	lines := make([]string, 0, len(allRows))
	for _, row := range allRows {
		lines = append(lines, fmt.Sprintf(
			"%-*s  %-*s  %-*s  %-*s  %s",
			branchWidth, row.Branch,
			headWidth, row.Head,
			stateWidth, row.State,
			modifiedWidth, row.Modified,
			row.Path,
		))
	}
	return lines
}

func worktreeSwitchCompletions(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	ctx, err := context.LoadWithProject()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	worktrees, err := worktree.List(ctx.Executor, ctx.Project.Path)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	seen := map[string]bool{}
	for _, wt := range worktrees {
		for _, candidate := range worktreeCompletionCandidates(wt) {
			if seen[candidate] || !strings.HasPrefix(candidate, toComplete) {
				continue
			}
			seen[candidate] = true
			completions = append(completions, fmt.Sprintf("%s\t%s", candidate, wt.Path))
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func worktreeCompletionCandidates(wt worktree.Worktree) []string {
	candidates := []string{}
	if wt.Branch != "" {
		candidates = append(candidates, wt.Branch)
	}
	base := filepath.Base(wt.Path)
	candidates = append(candidates, base)
	if name, ok := strings.CutPrefix(base, strings.Split(base, "--")[0]+"--"); ok {
		candidates = append(candidates, name)
	}
	return candidates
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
