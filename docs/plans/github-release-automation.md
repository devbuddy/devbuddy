# GitHub Release Automation Implementation Plan

> **For agentic workers:** execute task-by-task and keep each step verified before moving on.

**Goal:** make GitHub Actions the release authority so a coding agent can trigger a release without creating local release commits or tags.

**Architecture:** a tested Go helper computes the next release version from existing tags. A `release.yml` workflow runs from `workflow_dispatch`, validates the selected version, runs the same checks as CI, creates the tag, builds artifacts, publishes a GitHub Release with generated notes, and dispatches the Homebrew tap update.

**Tech Stack:** Go tests for version calculation, GitHub Actions, GitHub CLI for the local trigger wrapper.

---

### Task 1: Version Calculation Helper

**Files:**
- Create: `internal/release/version.go`
- Create: `internal/release/version_test.go`
- Create: `cmd/release-version/main.go`

Steps:
- Write tests for `minor`, `patch`, `rc`, and `custom` version selection.
- Verify the tests fail before implementation.
- Implement the helper and command.
- Verify the package tests pass.

### Task 2: Release Workflow

**Files:**
- Create: `.github/workflows/release.yml`
- Modify: `.github/workflows/tests.yml`

Steps:
- Add a manual release workflow with `kind` and `version` inputs.
- Run release validation and tests before tagging.
- Create an annotated tag from the workflow.
- Build binaries and checksums with `script/buildall`.
- Publish the GitHub Release with generated notes enabled.
- Keep the Homebrew dispatch after the GitHub Release publication.
- Remove the old tag-triggered release job from `tests.yml`.

### Task 3: Release Docs

**Files:**
- Modify: `CLAUDE.md`

Steps:
- Remove local release commands and wrappers.
- Document the agent-friendly `gh workflow run release.yml` trigger commands.
- Run formatting and focused tests.
