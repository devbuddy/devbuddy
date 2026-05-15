# Worktree UX Plan

## Summary

DevBuddy should make Git worktrees easy to create, maintain, and navigate for developers working across several branches, including AI-assisted work that often needs parallel checkouts.

The default disk layout is backward compatible:

```text
~/src/<platform>/<org>/<repo>
~/src/<platform>/<org>/<repo>--<worktree-name>
```

The existing clone path remains the main/default worktree. Additional worktrees are managed as sibling directories using the `repo--name` convention. This avoids nesting one Git checkout inside another and lets each worktree continue to behave like a normal DevBuddy project.

## Key UX Changes

- Add `bud wt` as the primary command group, with `bud worktree` as an alias.
- Keep `bud clone` unchanged by default. It still clones to the canonical repo path and jumps there.
- Make `bud cd` worktree-aware so fuzzy navigation can jump to either the canonical project or a managed sibling worktree.
- Add `bud wt list [query]` to show grouped worktrees with path, branch, HEAD short SHA, dirty state, and last modified time.
- Add `bud wt new <name> [branch]` to create `repo--<name>` from the current project repository.
- Add `bud wt cd [query]` to select and jump to a worktree, using the existing shell finalizer.
- Add `bud wt remove [query]` and `bud wt prune` for cleanup.

## Branch Conflict UX

When creating a worktree, DevBuddy should inspect `git worktree list --porcelain` before running `git worktree add`.

If the requested branch is already checked out in another worktree, DevBuddy should show a friendly conflict instead of surfacing Git's raw error:

- branch name
- existing worktree path
- dirty state when available
- suggested next actions

Interactive resolution choices:

- jump to the existing worktree
- create a new branch from the same commit, defaulting to `<branch>-<name>`
- create a detached worktree at the branch HEAD
- cancel

Non-interactive behavior should fail with a clear message and exact commands to resolve.

## Implementation Shape

- Add a `pkg/worktree` package for worktree discovery, parsing, path policy, and Git command construction.
- Parse `git worktree list --porcelain`; do not depend on ad hoc parsing of human-formatted output.
- Centralize path derivation so managed worktrees always live beside the canonical repo as `repo--<name>`.
- Keep Git execution isolated behind small functions so the code can later move behind the planned runtime environment abstraction.
- Avoid required global metadata in v1. Git worktree metadata and the filesystem layout are the source of truth.
- Extend project search to include both canonical projects and managed worktrees, while preserving existing fuzzy behavior.

## Test Plan

- Unit test porcelain parsing, including branch, detached HEAD, bare, and missing branch cases.
- Unit test sibling path derivation and worktree name sanitization.
- Unit test branch occupancy detection.
- Unit test project search with canonical repos and `repo--name` worktrees.
- Integration test `bud wt new feature-a` creates `repo--feature-a`.
- Integration test `bud wt cd feature-a` changes shell cwd through the existing finalizer.
- Integration test `bud cd feature-a` jumps directly to a managed worktree.
- Integration test branch conflict output for interactive and non-interactive use.
- Run shell finalizer coverage for both bash and zsh where integration coverage is added.

## Assumptions

- The sibling layout is the default: `repo--<worktree-name>`.
- Existing clones require no migration.
- The current repo path is the default/main worktree.
- `bud wt` is the preferred short command, and `bud worktree` is an alias.
- v1 focuses on local Git worktree creation, navigation, cleanup, and branch-conflict help. It does not include remote PR automation or AI-agent metadata.
