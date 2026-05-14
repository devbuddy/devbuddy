# Integration Test Improvements

## Context

DevBuddy has historically used a fully isolated Docker container with a real PTY for most integration tests. That gave strong end-to-end coverage, but two parts have been brittle:

- Prompt detection in the PTY/expect layer
- Writing test files into the container

The goal is not to preserve the current test form. The goal is to preserve and improve actual coverage while making normal development feedback faster and more reliable.

## Goals

- Keep enough real-environment coverage to prove the fast tests are not testing a fake world.
- Move tests that do not need PTY, Docker, or remote servers into cheaper harnesses.
- Avoid remote network calls in the default test suite.
- Keep an opt-in setting for allowing runtime tests to call real upstream servers.
- Move `github.com/devbuddy/expect` into this repository so shell-test infrastructure is easier to maintain.
- Make each step incremental and reversible.

## Test Layers

### Fast CLI Tests

These tests run the built `bud` binary directly with `os/exec`, temp directories, a temp `HOME`, and a temp `BUD_FINALIZER_FILE`. They should not use Docker, PTY, or remote network access.

They should cover:

- Basic commands: `bud`, `bud --help`, `bud --version`, `bud --debug-info`
- Manifest and task parse errors
- `bud inspect`
- `bud init`
- Custom command execution
- Custom task execution
- Filesystem effects from `bud up`
- Finalizer-file protocol for commands such as `bud cd`, `bud create`, and `bud up`

These tests prove the binary behavior, output, exit codes, project discovery, filesystem changes, and finalizer protocol. They do not prove parent-shell mutation.

### Focused Shell Tests

These tests keep real bash and zsh coverage, but only for behavior that actually depends on the user's shell.

They should cover:

- `eval "$(bud --shell-init)"` installs the wrapper and prompt hook
- The wrapper consumes `cd:` finalizers
- The wrapper consumes `setenv:` finalizers
- `__bud_prompt_command` preserves the previous command exit code
- Deferred initialization skips the first hook invocation
- Hook activation when entering a project
- Hook deactivation when leaving a project

Prefer a non-PTY shell runner if it is reliable:

- Start bash or zsh as a subprocess
- Send commands over stdin
- Append a unique completion sentinel after each command
- Parse exit status from the sentinel instead of detecting prompts

Keep one PTY smoke test only if we still need proof that interactive prompt-hook behavior works in a real terminal session.

### Docker Runtime Tests

Docker should remain for Linux runtime behavior that depends on the target environment, installed tools, or distro assumptions.

The default Docker suite should avoid remote servers. It should use preseeded or mounted fixtures for:

- Go archives
- Node archives
- Python or pyenv artifacts
- Python wheels
- Go modules
- npm packages

These tests should prove DevBuddy's orchestration and activation behavior, not the availability of upstream services.

The same runtime tests should be able to run in a real-network mode. In that mode, fixture redirects are disabled and DevBuddy calls the normal upstream destinations.

Recommended control:

- `BUD_TEST_REAL_NETWORK=1`

This mode can verify real downloads and package installation against actual Go, Node, PyPI, npm, or other upstream services. It should not run by default on every pull request.

## Incremental Plan

### 1. Label Existing Coverage

Introduce clear test categories before moving tests:

- `fast`
- `shell`
- `docker`

This can start as package names, build tags, CI job names, or helper names. The important part is making the intended cost and environment visible.

### 2. Move `devbuddy/expect` Into This Repository

Make this a mechanical change first:

- Copy the expect package into this repository, likely under `tests/internal/expect` or `internal/expect`
- Update imports from `github.com/devbuddy/expect`
- Do not change harness behavior in the same step
- Run the current bash and zsh integration tests

After this, prompt handling, debug output, and timeout behavior can be improved locally without coordinating another repository.

### 3. Replace Container File Uploads With Mounted Workspaces

Keep the current Docker and PTY strategy for this step, but stop writing files through the interactive shell.

Use a host-side temp directory mounted into the container:

- Tests create and mutate files with normal Go filesystem APIs
- Docker mounts the temp directory into a stable container path
- Container commands run against the mounted workspace

This removes one brittle area without changing shell behavior.

### 4. Add A Fast Black-Box CLI Harness

Build `bud` once in `TestMain`, then run it directly from temp project directories.

The harness should provide:

- `Run` with cwd, env, expected exit code, and timeout
- `Write`, `Read`, `Exists`, and directory assertions using host filesystem APIs
- A per-command `BUD_FINALIZER_FILE`
- Helpers to inspect finalizer contents

Start by migrating tests that obviously do not need parent-shell mutation.

Good first candidates:

- Command help/version/debug-info
- Inspect
- Task parse errors
- Custom commands
- Custom tasks
- Init

### 5. Shrink The PTY Suite

Once fast CLI coverage exists, remove Docker/PTY usage from tests that do not need it.

The remaining PTY or shell tests should be small and intentional:

- Shell initialization
- Prompt hook behavior
- Parent-shell environment mutation
- Finalizer consumption by the shell wrapper

### 6. Try A Non-PTY Shell Runner

Build a subprocess-based shell harness and compare it against the remaining shell cases.

If it handles bash and zsh reliably, migrate most shell tests to it. Keep a single PTY smoke test if needed for confidence in interactive prompt behavior.

### 7. Split Runtime Install Tests

Move Go, Node, Python, pip, pipfile, and python develop coverage into Docker runtime tests with offline fixtures by default.

Then add `BUD_TEST_REAL_NETWORK=1` as a mode that lets those same tests call the normal upstream destinations.

### 8. Update CI

Normal pull request CI should run:

- Package unit tests
- Fast CLI tests
- Focused bash shell tests
- Focused zsh shell tests
- Offline Docker runtime tests if they are reasonably fast

Scheduled or manual CI should run:

- Full Docker runtime tests
- Docker runtime tests with `BUD_TEST_REAL_NETWORK=1`

## Success Criteria

- Most pull request feedback no longer depends on Docker, PTY, or remote servers.
- The remaining shell tests directly cover shell integration behavior.
- Runtime installation tests are deterministic by default.
- Real upstream integration remains available through an explicit runtime-test setting.
- The local expect package can be modified alongside the test harness.
- Coverage is easier to explain from the test name, package, or CI job.
