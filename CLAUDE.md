# DevBuddy - Claude Instructions

## Overview

DevBuddy (`bud`) is a development environment manager. It reads a `dev.yml` file in a project root and automates environment setup (language versions, virtualenvs, system packages) and provides custom project commands. It integrates deeply with the user's shell (bash/zsh) via a shell hook that automatically activates/deactivates environments as you `cd` between projects.

### Purpose & scope

DevBuddy is an open-source equivalent of Shopify's internal `dev` tool. The goal is to keep a comparable scope to `dev`: manage the **whole developer environment**, not just language versions. That includes:

- Language toolchains and project virtualenvs (Go, Python, Node, Ruby, ...)
- System packages (Homebrew, apt) and platform-level prerequisites (e.g. accepting the Xcode license after an OS upgrade)
- Project ergonomics: cloning, `cd`ing between repos, opening project URLs, project-local commands
- Future: worktree management and other repo-level workflows

This broader scope is intentional and is what differentiates DevBuddy from pure version managers.

### Related projects (reference material)

Two well-maintained projects overlap with DevBuddy on the version-manager axis and are good references when designing language tasks:

- **asdf** — https://github.com/asdf-vm/asdf — pluggable, single-tool version manager; `.tool-versions` file; rich plugin ecosystem.
- **mise** — https://github.com/jdx/mise — Rust-based asdf-compatible version manager; also handles env vars and project tasks; very active.

When adding or revising language support, check how asdf/mise handle the same case (version-file conventions, plugin/install flow, shimming vs PATH activation, lockfile watching). DevBuddy's *scope* is wider than theirs, but its *version-manager surface* should not be worse than theirs for the languages it supports.

## Roadmap

See `ROADMAP.md` for the product and technical direction. Consult it before making architectural changes to ensure they align with the desired direction. Key points:
- Moving toward a **runtime environment abstraction** (command execution, filesystem, terminal I/O, env vars, state store) with real and testing implementations
- Near-term: refactor `pkg/executor/` from interface+builder to a plain struct representing a command request
- Changes should move toward testability without Docker/PTY and away from direct OS calls in business logic

## Planning

Write implementation plans in `docs/plans/`. Keep plans concrete enough to execute incrementally, with file paths, validation commands, and PR-sized steps.

## Project Structure

```
cmd/bud/main.go          # Entry point. Version set via build flags (-ldflags)
pkg/
  cmd/                    # Cobra CLI commands (root, up, cd, clone, open, create, init, inspect, upgrade, commands)
  config/                 # Global config (debug mode via BUD_DEBUG env var)
  context/                # Runtime context: loads project, env, UI
  env/                    # Environment variable mutation tracking (set, prepend, compare)
  executor/               # Command execution helpers
  autoenv/                # Feature activation/deactivation state machine
    feature_info.go       # FeatureInfo (name+param), FeatureSet
    state.go              # Persisted state of active features (JSON in BUD_AUTO_ENV_FEATURES)
    runner.go             # Sync() - activates/deactivates features by diffing desired vs current
  hook/                   # Shell hook logic (called on every prompt via `bud --shell-hook`)
  integration/            # Shell integration scripts (embedded via //go:embed)
    common.sh             # bud() wrapper, __bud_prompt_command, finalizers
    bash.sh               # PROMPT_COMMAND hook
    zsh.sh                # precmd_functions hook
    integration.go        # Print(), CompletionScriptProvider, AddFinalizerCd
  manifest/               # dev.yml parsing
  project/                # Project discovery (walks up to find dev.yml)
  tasks/                  # Task implementations
    api/                  # TaskAction interface, Condition interface (file checksum tracking)
    taskengine/           # Task execution engine
    golang.go             # Go version management (installs Go, sets GOROOT/GOPATH/GO111MODULE)
    python.go             # Python via pyenv
    python_develop.go     # pip install -e . (tracks setup.py/pyproject.toml changes)
    pip.go                # pip install -r requirements.txt
    pipfile.go            # Pipfile support
    node.go               # Node.js version management
    homebrew.go           # Homebrew package installation
    apt.go                # APT package installation
    custom.go             # Custom tasks (met?/meet pattern)
    env.go                # Environment variable task
    envfile.go            # .env file loading
    golang_dep.go         # Go dep (legacy)
  helpers/                # Shared helpers
    debug/                # Debug info collector
    downloader.go         # HTTP file downloader
    git.go                # Git operations
    github.go             # GitHub URL parsing
    golang.go             # Go installation helper
    homebrew.go           # Homebrew detection/installation
    node.go               # Node.js installation
    pyenv.go              # Pyenv management
    virtualenv.go         # Python virtualenv creation
    store/                # Key-value store (.devbuddy/ dir in project)
    open/                 # URL opener (browser)
    osidentity/           # OS detection (linux/darwin)
    projectmetadata/      # Project source path conventions
    fixtures/             # Test fixtures (VCR cassettes)
  termui/                 # Terminal UI (colored output, spinners, panda emoji)
  test/                   # Test utilities
  utils/                  # File utils, checksums, path helpers
script/
  test                    # Run unit tests (./pkg/...)
  lint                    # Run golangci-lint
  buildall                # Cross-compile for all platforms
  release                 # Dispatch the GitHub release workflow
  install-dev.sh          # Build and install to GOPATH/bin
tests/                    # Integration tests (Docker-based)
  context/                # TestContext: orchestrates Docker container with shell
  main_test.go            # TestMain: builds Linux binary, loads config
  helper_test.go          # CreateContext, CreateProject, output assertions
  cmd_*_test.go           # Command integration tests
  task_*_test.go          # Task integration tests
```

## Key Architectural Concepts

### Shell Integration
The user adds `eval "$(bud --shell-init)"` to their shell config. This:
1. Defines a `bud()` shell function that wraps the real binary
2. The wrapper handles "finalizers" (cd, setenv) via a temp file (`BUD_FINALIZER_FILE`)
3. Installs `__bud_prompt_command` as PROMPT_COMMAND (bash) or precmd (zsh)
4. On every prompt, `bud --shell-hook` runs and outputs shell commands to stdout
5. The hook output is `eval`'d to export/unset env vars, activate virtualenvs, etc.

Shell hook guardrails:
- `__bud_prompt_command` must preserve the previous command exit code (`$?`) so prompt-hook execution never clobbers user-visible status.
- Changes to shell hooks should include an integration regression test under `tests/` (for both bash and zsh jobs).

### AutoEnv (Feature Activation)
- Each task can declare a "feature" (e.g., `python=3.6.5`, `golang=1.21`)
- The hook tracks active features in `BUD_AUTO_ENV_FEATURES` env var (JSON state)
- `autoenv.Sync()` diffs desired vs active features, activating/deactivating as needed
- Features mutate the environment (PATH, GOROOT, VIRTUAL_ENV, etc.)
- When leaving a project, features are deactivated and env vars restored

### Task System
- Tasks implement `TaskAction` interface: `Description()`, `Needed()`, `Run()`, `Feature()`
- `Needed()` returns whether the task should run (idempotent check)
- `Condition` interface supports file-checksum-based change detection
- Task state stored in `.devbuddy/` directory within the project

### Finalizer Mechanism
Some commands (like `bud cd`) need to change the shell's working directory. Since a subprocess can't change the parent shell's cwd, DevBuddy writes a "finalizer" to `BUD_FINALIZER_FILE`, and the shell wrapper function processes it after the command exits.

## Development

### Build & Run
```bash
bud install              # Build and install to GOPATH/bin
go build -o bud ./cmd/bud  # Quick local build
```

### Testing
```bash
script/test              # Unit tests only (./pkg/...)
script/lint              # golangci-lint
```

### Integration Tests
Don't use `bud integration` (requires a PTY). Run the go test command directly:
```bash
export TEST_DOCKER_IMAGE="ghcr.io/devbuddy/docker-testing:sha-7fd13f4"
TEST_SHELL=bash go test -v -count=1 ./tests
TEST_SHELL=zsh go test -v -count=1 ./tests
```
- Takes ~1 minute, 34 pass / 5 skipped (Node platform not available, env var fixmes)
- Use `-count=1` to bypass test cache
- Requires Docker running locally

### Integration Test Architecture
Integration tests run inside a Docker container (`ghcr.io/devbuddy/docker-testing`):
- `TestMain()` cross-compiles a Linux binary, mounts it into the container
- `TestContext` uses `tests/internal/expect` (PTY-based shell automation)
- Tests create projects with dev.yml, run `bud` commands, assert output
- Controlled by env vars: `TEST_SHELL` (bash/zsh), `TEST_DOCKER_IMAGE`
- Docker image: Ubuntu 20.04 with pyenv, Python 3.9, build tools, git, zsh

### CI/CD
GitHub Actions:
- `.github/workflows/tests.yml`: golangci-lint, unit tests, bash integration, zsh integration, task integration
- `.github/workflows/release.yml`: manually triggered release authority
- Go 1.26, golangci-lint v2.9.0

### Release Process
```bash
script/release minor              # Minor bump (0.16.0 -> 0.17.0)
script/release patch              # Patch bump (0.16.0 -> 0.16.1)
script/release rc                 # Next release candidate
script/release custom v0.17.0     # Explicit version
script/release --dry-run patch    # Print the workflow command
```
`script/release` only dispatches the GitHub workflow. The workflow computes or validates the version, runs release validation, creates the annotated tag, builds binaries via `script/buildall`, creates a GitHub Release with generated notes and attached binaries, and triggers the Homebrew tap update for stable releases.

Coding agents can trigger releases directly with:
```bash
gh workflow run release.yml --repo devbuddy/devbuddy --ref main -f kind=patch
```

The Homebrew tap update requires a valid `ACCESS_TOKEN` secret.

### Distribution
- macOS: Homebrew (`devbuddy/homebrew-devbuddy`, auto-triggered on release)
- Linux: GitHub Releases (direct binary download)
- Platforms: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64

## Dependencies
- `github.com/spf13/cobra` - CLI framework
- `tests/internal/expect` - PTY-based shell automation (moved from `github.com/devbuddy/expect`)
- `github.com/joho/godotenv` - .env file parsing
- `github.com/logrusorgru/aurora` - Terminal colors
- `github.com/sahilm/fuzzy` - Fuzzy matching (for `bud cd`)
- `github.com/mitchellh/go-ps` - Process detection (shell identification)
- `github.com/dnaeon/go-vcr` - HTTP recording for tests
- `github.com/goccy/go-yaml` - YAML parsing

## Conventions
- Version is set at build time via `-ldflags` (not in source)
- No changelog maintained
- Panda emoji in UI output
- `.devbuddy/` directory stores per-project state (checksums, etc.)
- `BUD_DEBUG=1` enables debug logging in both Go code and shell hooks
- When creating or updating a pull request description for an issue-driven change, include a closing reference in this exact format: `Fixes: #<ISSUE-NUMBER>`
