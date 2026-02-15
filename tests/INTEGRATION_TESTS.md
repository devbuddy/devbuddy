# Integration Tests

## Goals

The integration tests verify that DevBuddy works correctly as an end-user would experience it:
running inside a real shell, with real environment mutations, real task execution, and real
shell integration (hooks, finalizers, autoenv).

This is necessary because DevBuddy's core value is shell integration — it mutates the user's
shell environment, and unit tests cannot verify that. The integration tests prove that:

1. **Shell integration works**: `eval "$(bud --shell-init)"` installs the wrapper and hooks correctly
2. **Shell hook works**: entering/leaving a project activates/deactivates features (env vars, PATH)
3. **Finalizers work**: `bud cd` changes the shell's working directory
4. **Tasks execute correctly**: `bud up` runs tasks, respects conditions, reports errors
5. **Custom commands work**: project commands from dev.yml execute with the right environment
6. **Both bash and zsh work**: the shell integration must work identically in both shells

## Requirements

### Shell coverage

- **Bash and zsh are both essential.** Users use both, and the shell integration code differs
  between them (PROMPT_COMMAND vs precmd_functions, different quoting edge cases).
- Every test must pass in both shells.

### Platform coverage

- **Linux**: primary CI platform, tested via Docker containers.
- **macOS**: needed for macOS-specific tasks (homebrew). Currently not covered by integration
  tests. This is a gap.

### Isolation

- Each test must start from a clean, predictable state.
- Tests must not affect each other. A failing test must not cause subsequent tests to fail.
- The test environment must not depend on the host machine's state beyond Docker being available.

### No network dependency

- Integration tests must not depend on external network calls (downloading Go, Python, Node, etc.).
- External downloads are a major source of flakiness: slow, rate-limited, or unavailable.
- **Strategy**: use a local HTTP server serving dummy payloads, or pre-bake tools into the
  Docker image, or force the HTTP client to target a local server.
- The tests should prove that DevBuddy correctly orchestrates the installation — not that the
  upstream download servers are available.

### Robustness

- Tests must be deterministic. Flaky tests erode trust and slow down development.
- Prompt detection must be reliable. The PTY/expect layer is the most fragile part.
- Timeouts must be generous enough for slow CI machines but tight enough to catch hangs.
- CI and local-dev must behave the same way. A test that passes in CI but fails locally
  (or vice versa) is a bug in the test infrastructure.

### Speed

- Fast feedback is important. Tests that take minutes to run get skipped.
- Container startup/teardown overhead should be minimized.
- Tests that don't need network or heavy setup should run fast.

## Current Implementation

### Architecture

```
Host machine (macOS or Linux CI)
  │
  ├─ TestMain (tests/main_test.go)
  │   ├─ Loads config from env vars (TEST_SHELL, TEST_DOCKER_IMAGE)
  │   └─ Cross-compiles a static Linux binary: GOOS=linux CGO_ENABLED=0
  │
  ├─ Per test function:
  │   ├─ CreateContext() → starts a Docker container with the compiled binary mounted
  │   │   ├─ docker run -ti --rm -v bud:/usr/local/bin/bud -e PS1=... IMAGE SHELL
  │   │   ├─ Wraps with expect.ShellExpect (prompt-based command/response)
  │   │   ├─ Disables echo: stty -echo
  │   │   └─ Verifies IN_DOCKER=yes
  │   │
  │   ├─ CreateContextAndInit() → also runs: eval "$(bud --shell-init)"
  │   │
  │   ├─ CreateProject() → creates a temp directory with a dev.yml inside the container
  │   │
  │   └─ c.Run(t, "bud up") → sends command, waits for prompt, checks exit code
  │
  └─ Container cleanup via t.Cleanup() → sends SIGKILL to Docker process
```

### Key components

**tests/main_test.go** — TestMain builds the Linux binary once before all tests run.

**tests/context/** — TestContext wraps the Docker + PTY interaction:
- `context.go`: New(), Run(), Write(), Cat(), Cd(), GetEnv(), Close()
- `config.go`: loads TEST_SHELL (default "bash"), TEST_DOCKER_IMAGE (required)
- `options.go`: Timeout() (default 5s), ExitCode() (default 0)
- `strip_ansi.go`: strips ANSI escape codes from command output

**tests/helper_test.go** — Test helpers:
- `CreateContext(t)` / `CreateContextAndInit(t)`: container lifecycle
- `CreateProject(t, c, devYmlLines...)`: creates a project dir with dev.yml
- `OutputContains(t, lines, ...)` / `OutputEqual(t, lines, ...)`: output assertions

**github.com/devbuddy/expect** — PTY-based shell automation library (custom fork of etcd's):
- `ExpectProcess`: low-level PTY process (start, send, read lines, stop)
- `ShellExpect`: sends commands, waits for prompt string, returns output
- Prompt detection: reads lines until accumulated output ends with the prompt string

**ghcr.io/devbuddy/docker-testing** — Pre-built Docker image:
- Base: Ubuntu 20.04
- Includes: bash, zsh, git, curl, build tools, pyenv, Python 3.9.0
- User: `tester` (non-root), home: `/home/tester`
- Published to GitHub Container Registry, pinned by SHA in CI

### How a test runs a command

1. `c.Run(t, "bud up")` calls the internal `run()` method
2. `run()` calls `c.shell.Run(cmd)` in a goroutine (for timeout)
3. `ShellExpect.Run()` sends `cmd + "\n"` to the PTY
4. It reads lines until the accumulated output ends with the prompt string
5. Back in `run()`, it sends `echo $?` to get the exit code
6. It parses the exit code and compares with the expected one
7. Output lines are stripped of ANSI codes and returned

### Test patterns

**Simple command test** (no shell-init needed):
```go
func Test_Cmd_Help(t *testing.T) {
    c := CreateContext(t)
    lines := c.Run(t, "bud")
    OutputContains(t, lines, "Usage:", "DevBuddy Commands:")
}
```

**Task test** (needs shell-init for hooks):
```go
func Test_Task_Custom(t *testing.T) {
    c := CreateContextAndInit(t)
    p := CreateProject(t, c, `up:\n- custom:\n    met?: test -e sentinel\n    meet: echo A > sentinel`)
    c.Cd(t, p.Path)
    c.Run(t, "bud up")
    content := c.Cat(t, "sentinel")
    require.Equal(t, "A", content)
}
```

**Shared container via subtests** (avoids redundant downloads):
```go
func Test_Task_Go(t *testing.T) {
    c := CreateContextAndInit(t)
    t.Run("installs_and_runs_go_modules", func(t *testing.T) { /* reuses container */ })
    t.Run("modules_false_is_rejected", func(t *testing.T) { /* reuses container */ })
}
```

### Environment variables

| Variable | Required | Default | Purpose |
|----------|----------|---------|---------|
| `TEST_SHELL` | no | `bash` | Shell to test (bash or zsh) |
| `TEST_DOCKER_IMAGE` | **yes** | — | Docker image for the container |

### CI configuration

`.github/workflows/tests.yml` runs four parallel jobs:
1. `golangci-lint` — linting
2. `go test ./pkg/...` — unit tests
3. `TEST_SHELL=bash go test -v ./tests` — bash integration
4. `TEST_SHELL=zsh go test -v ./tests` — zsh integration

The release job depends on all four passing.

## Known Issues

### Prompt detection is fragile

The expect library detects command completion by matching the prompt string
(`\nPROMPTPROMPTPROMPTPROMPT\n`) at the end of accumulated output. This breaks when:
- A command's output happens to contain the prompt string
- The PTY splits the prompt across multiple reads
- Terminal initialization output interferes (note: `c.Init()` is commented out in context.go)

### One container per test function

Each call to `CreateContext()` starts a new Docker container. This gives good isolation but is
slow. Container startup typically takes 1-3 seconds. Tests that share a container via subtests
(like `Test_Task_Go`) work around this but sacrifice isolation — a failing subtest can leave
the container in a dirty state.

### Timeout handling

The default timeout is 5 seconds. Tasks that download tools (Go, Python, Node) need
`context.Timeout(2*time.Minute)`. If the timeout is too short, the test fails with a
confusing "timed out" error rather than the actual command output. There is no way to see
what the command was doing when it timed out.

### Docker image version mismatch

`dev.yml` pins `TEST_DOCKER_IMAGE: ghcr.io/devbuddy/docker-testing:sha-7fd13f4` but CI uses
`ghcr.io/devbuddy/docker-testing:sha-f11e362`. This means local `bud integration` and CI may
use different images.

### No macOS integration tests

Tasks like `homebrew` that are macOS-specific cannot be tested in the Docker-based setup
(which runs Linux). There is currently no way to run integration tests against macOS.

### Network-dependent tests are slow and flaky

Tests for Go, Python, and Node tasks download real distributions from the internet.
This makes them slow (~minutes) and vulnerable to network issues or rate limiting.

### Skipped tests

`Test_Task_Custom_With_Env_At_First_Run` is skipped with `t.Skip("Fixme: env vars not set
before tasks?")`. This indicates a known bug where dev.yml `env:` vars are not available
during the first `bud up` run (they're only set by the shell hook, which runs on prompt, not
during `bud up`).
