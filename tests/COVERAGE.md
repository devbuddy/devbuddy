# Integration Test Coverage Labels

This file labels the existing integration tests before the harness is split. The
labels describe the coverage each test should keep as the suite moves away from
the current Docker + PTY default.

All tests in `tests/` currently run through the Docker + PTY harness. The
`target label` column is the intended destination or retained coverage layer.

## Labels

| Label | Meaning |
| --- | --- |
| `fast` | Black-box `bud` binary behavior that can run with `os/exec`, temp directories, and a temp `BUD_FINALIZER_FILE`. No Docker, PTY, or remote network should be required. |
| `shell` | Behavior that must prove real shell integration: wrapper setup, prompt hook behavior, parent-shell cwd/env mutation, or bash/zsh-specific behavior. |
| `docker` | Runtime behavior that should continue to run in a Linux container because it depends on distro assumptions, installed tools, or language runtime installation. |

`docker` tests should avoid remote servers by default. When
`BUD_TEST_REAL_NETWORK=1` is set, the same runtime tests may call real upstream
destinations instead of fixture or preseeded artifacts.

## Current Test Inventory

| File | Current focus | Target label |
| --- | --- | --- |
| `tests/fast/cmd_test.go` | Help, version, and debug-info command output. | `fast` |
| `tests/fast/cmd_inspect_test.go` | Manifest inspection and project-not-found error output. | `fast` |
| `tests/fast/cmd_open_test.go` | Custom project links, fuzzy link matching, default link behavior, and missing-link-name errors. | `fast` |
| `tests/fast/task_error_test.go` | Manifest and task parse errors from `bud up`. | `fast` |
| `tests/fast/cmd_custom_test.go` | Custom command execution, output, env injection, stdin, exit-code normalization, project-root cwd. | `fast` |
| `tests/fast/task_custom_test.go` | Custom task conditions, run/skip behavior, project-root cwd, task env behavior. | `fast` |
| `tests/shell/cmd_init_test.go` | Manifest template creation plus shell-hook activation after `bud init`. | `fast`, `shell` |
| `tests/shell/cmd_cd_test.go` | Project matching plus wrapper consumption of `cd:` finalizers. | `fast`, `shell` |
| `tests/shell/cmd_create_test.go` | Project creation, manifest template creation, `cd:` finalizer, activation after create. | `fast`, `shell` |
| `tests/shell/cmd_hook_test.go` | Prompt hook exit-code preservation, deferred init, immediate activation. | `shell` |
| `tests/shell/global_env_test.go` | Parent-shell env activation and cleanup when leaving a project. | `shell` |
| `tests/shell/env_envfile_test.go` | `.env` feature activation, cleanup, in-process availability, and change detection. | `shell` |
| `tests/docker/task_go_test.go` | Go runtime installation, activation, module command behavior, and Go task validation. | `docker` |
| `tests/docker/task_node_test.go` | Node runtime installation, activation, and npm install behavior. | `docker` |
| `tests/docker/task_python_test.go` | Python runtime activation and virtualenv behavior. | `docker` |
| `tests/docker/task_pip_test.go` | Pip install into activated Python environment. | `docker` |
| `tests/docker/task_pipfile_test.go` | Pipfile install into activated Python environment. | `docker` |
| `tests/docker/task_python_develop_test.go` | Editable Python package install, update detection, and extras handling. | `docker` |

## Migration Notes

- Tests with multiple labels should be split when migrated. For example,
  `cmd_init_test.go` should keep manifest creation coverage in `fast` and keep
  activation coverage in `shell`.
- `cmd_custom_test.go` and `task_custom_test.go` are good first migration
  candidates because most assertions are filesystem, output, and exit-code
  behavior rather than parent-shell mutation.
- Runtime install tests remain `docker` even if a future fast harness verifies
  their parsing or command-construction behavior elsewhere.
