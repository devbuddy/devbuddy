# DevBuddy Roadmap

Product and technical direction for the project. Consult this when making changes to ensure they align with the desired direction.

## Product Direction

DevBuddy automates development environment setup. The core value is: `cd` into a project, everything activates; run `bud up`, everything installs. It should work reliably across shells, platforms, and terminal environments (interactive, CI, editors, piped).

## Technical Direction

### Runtime Environment Abstraction

The long-term goal is an **environment service** — a single interface representing the execution environment that all code operates through. This covers:

- **Command execution** — running subprocesses
- **Filesystem** — reading/writing files, checking paths
- **Terminal I/O** — stdout, stderr, stdin, isatty detection
- **Environment variables** — reading, setting, tracking mutations
- **State store** — per-project state (`.devbuddy/` directory)
- **Logging / UI** — debug output, user-facing messages

Production code uses a **real implementation** that talks to the actual OS. Tests use a **testing implementation** that can be inspected and controlled without real subprocesses, filesystem writes, or terminal access.

This makes the codebase testable without Docker containers or PTY automation, and eliminates the class of bugs where code assumes a specific execution context (interactive terminal, specific OS, writable filesystem).

### Near-term: Executor Refactor

The first step toward the environment abstraction is refactoring `pkg/executor/`:

**Current state:** an `Executor` interface with 11 methods, a fluent builder pattern, and 3 execution modes (output-filter, PTY, passthrough). Tests need type assertions to access internal fields. Several interface methods are unused.

**Target state:** a plain struct representing a command execution request:

```go
type ExecTask struct {
    Command      string
    Args         []string
    Shell        bool       // if true, run via sh -c
    Cwd          string
    Env          []string
    Passthrough  bool       // inherit stdin/stdout/stderr directly
    OutputPrefix string
}

func (t ExecTask) Run() *Result { ... }
func (t ExecTask) Capture() *Result { ... }
```

Benefits:
- Tests inspect struct fields directly — no mocking, no type assertions
- Adding a field doesn't break an interface
- No ambiguity about configuration order or repeated calls
- Natural stepping stone to injecting a `RunFunc` for testing

**Drop:** `SetPTY` (unused), `AddOutputFilter` (unused), `SetEnvVar` (unused), the `Executor` interface itself.

### Future Steps (not yet planned)

Once the executor is a plain struct, introduce a `Runtime` or `Env` service:

1. Extract terminal I/O from `termui` into the service
2. Move env var tracking (`pkg/env/`) behind the service
3. Move state store (`pkg/helpers/store/`) behind the service
4. Move filesystem operations behind the service
5. Inject the service into task context, replacing direct OS calls

Each step is independently valuable and shippable.
