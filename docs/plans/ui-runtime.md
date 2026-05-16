# UI Runtime Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the ad hoc `termui`/prompt handling with a testable UI runtime seam that centralizes terminal output, ANSI policy, prompt behavior, and recorded UI events.

**Architecture:** Introduce a small UI facade that records intent-level events and delegates final text formatting to renderers. Production code uses a terminal renderer and real prompts; tests use a plain renderer and fake prompts so command behavior can be asserted without ANSI stripping, Docker, or PTY automation unless the test is explicitly about shell integration.

**Tech Stack:** Go, Cobra, existing `pkg/context`, existing `pkg/termui` compatibility surface, `github.com/charmbracelet/huh` for prompts behind an internal interface, existing `script/test` and shell integration tests.

---

## Constraints

- Keep behavior and visible output stable unless a task explicitly says otherwise.
- Do not expose Charmbracelet types outside the prompt adapter. Commands and tasks should depend on DevBuddy-owned request/result structs.
- Prefer plain unit tests for UI decisions. Keep PTY tests only for shell-finalizer and real prompt smoke coverage.
- Use `docs/plans/` for follow-up plans.
- Do not remove `termui` in the first PR; migrate behind a compatibility layer first.

## File Map

- Create `pkg/ui/event.go`: event kinds, fields, and the recorder-friendly event model.
- Create `pkg/ui/renderer.go`: renderer interface, plain renderer, terminal renderer, color policy.
- Create `pkg/ui/ui.go`: production/testing constructors, writers, emitted event storage, and high-level UI methods.
- Create `pkg/ui/prompts.go`: prompt request/result types and prompt interface.
- Create `pkg/ui/huh_prompts.go`: production prompt adapter using `huh`, added after the prompt interface exists.
- Create `pkg/ui/fake_prompts.go`: testing prompt implementation.
- Modify `pkg/termui/*`: make `termui.UI` a compatibility wrapper over `pkg/ui.UI`.
- Modify `pkg/context/context.go`: construct the new UI once and pass it through the existing context field.
- Modify prompt call sites: `pkg/cmd/worktree.go`, `pkg/cmd/init.go`, `pkg/cmd/clone.go`.
- Modify tests that currently inspect rendered text where event assertions are clearer.
- Later cleanup: remove direct `aurora` usage from `pkg/termui`, `pkg/cmd/root.go`, and `pkg/helpers/open/open.go`.

## Stack Overview

1. **PR 1: UI event core and termui compatibility.** Add `pkg/ui`, keep output stable, and add event assertions for existing unit tests.
2. **PR 2: ANSI policy centralization.** Move color rendering into `pkg/ui.Renderer`, default tests to plain output, and reduce ANSI stripping.
3. **PR 3: Prompt interface and fake prompts.** Add DevBuddy-owned prompt interfaces and migrate worktree selection logic to use prompt requests in unit-testable code.
4. **PR 4: Replace promptui with huh adapter.** Use `huh` only behind `pkg/ui.Prompts`, verify ESC cancel and keep one PTY smoke test.
5. **PR 5: Migrate remaining prompts and remove promptui.** Move `init` template selection and `clone` confirmation, then drop `promptui`.

Each PR should be based on the previous branch and opened as a stacked PR. After each push, monitor GitHub Actions until CI passes before starting the next PR.

---

## PR 1: UI Event Core And termui Compatibility

**Branch:** `ui-runtime-core`

**Files:**
- Create: `pkg/ui/event.go`
- Create: `pkg/ui/renderer.go`
- Create: `pkg/ui/ui.go`
- Create: `pkg/ui/ui_test.go`
- Modify: `pkg/termui/ui.go`
- Modify: `pkg/termui/task.go`
- Modify: `pkg/termui/hook.go`
- Modify: `pkg/context/context.go`
- Modify: focused tests under `pkg/autoenv`, `pkg/tasks/taskengine`, and `pkg/context`

- [ ] **Step 1: Create the event model**

Create `pkg/ui/event.go`:

```go
package ui

type Kind string

const (
	KindDebug             Kind = "debug"
	KindWarning           Kind = "warning"
	KindCommandHeader     Kind = "command_header"
	KindCommandRun        Kind = "command_run"
	KindCommandActed      Kind = "command_acted"
	KindProjectExists     Kind = "project_exists"
	KindJumpProject       Kind = "jump_project"
	KindTaskHeader        Kind = "task_header"
	KindTaskCommand       Kind = "task_command"
	KindTaskShell         Kind = "task_shell"
	KindTaskActed         Kind = "task_acted"
	KindTaskAlreadyOK     Kind = "task_already_ok"
	KindTaskError         Kind = "task_error"
	KindTaskWarning       Kind = "task_warning"
	KindTaskActionHeader  Kind = "task_action_header"
	KindActionHeader      Kind = "action_header"
	KindActionNotice      Kind = "action_notice"
	KindActionDone        Kind = "action_done"
	KindHookActivated     Kind = "hook_activated"
	KindHookFeatureFailed Kind = "hook_feature_failed"
	KindHookDevYMLChanged Kind = "hook_devyml_changed"
)

type Field struct {
	Name  string
	Value string
}

type Event struct {
	Kind   Kind
	Text   string
	Fields []Field
}

func F(name, value string) Field {
	return Field{Name: name, Value: value}
}
```

- [ ] **Step 2: Add renderers with stable current output**

Create `pkg/ui/renderer.go` with a `Renderer` interface and two implementations:

```go
package ui

import (
	"fmt"
	"strings"

	color "github.com/logrusorgru/aurora"
)

type Renderer interface {
	Render(Event) string
}

type PlainRenderer struct{}
type TerminalRenderer struct{}

func (PlainRenderer) Render(event Event) string {
	return renderPlain(event)
}

func (TerminalRenderer) Render(event Event) string {
	return renderTerminal(event)
}

func renderPlain(event Event) string {
	switch event.Kind {
	case KindDebug:
		return fmt.Sprintf("BUD_DEBUG: %s\n", event.Text)
	case KindWarning:
		return fmt.Sprintf("WARNING: %s\n", event.Text)
	case KindCommandHeader:
		return fmt.Sprintf("🐼  running %s\n", event.Text)
	case KindCommandRun:
		return fmt.Sprintf("%s %s\n", event.Text, field(event, "args"))
	case KindCommandActed:
		return "  Done!\n"
	case KindProjectExists:
		return "🐼  project already exists locally\n"
	case KindJumpProject:
		return fmt.Sprintf("🐼  jumping to %s\n", event.Text)
	case KindTaskHeader:
		return renderTaskHeaderPlain(event)
	case KindTaskCommand:
		return fmt.Sprintf("  Running: %s %s\n", event.Text, field(event, "args"))
	case KindTaskShell:
		return fmt.Sprintf("  Running: %s\n", event.Text)
	case KindTaskActed:
		return "  Done!\n"
	case KindTaskAlreadyOK:
		return "  Already OK!\n"
	case KindTaskError:
		return fmt.Sprintf("  %s\n", event.Text)
	case KindTaskWarning:
		return fmt.Sprintf("  Warning: %s\n", event.Text)
	case KindTaskActionHeader:
		return fmt.Sprintf("  ▪︎%s\n", event.Text)
	case KindActionHeader:
		return fmt.Sprintf("🐼  %s\n", event.Text)
	case KindActionNotice:
		return fmt.Sprintf("⚠️   %s\n", event.Text)
	case KindActionDone:
		return "✅  Done!\n"
	case KindHookActivated:
		return fmt.Sprintf("🐼  activated: %s\n", event.Text)
	case KindHookFeatureFailed:
		param := field(event, "param")
		if param != "" {
			param = fmt.Sprintf(" (%s)", param)
		}
		return fmt.Sprintf("🐼  failed to activate %s. Try running 'bud up' first!%s\n", event.Text, param)
	case KindHookDevYMLChanged:
		return "🐼  dev.yml changed, run `bud up` to apply\n"
	default:
		return event.Text + "\n"
	}
}

func renderTerminal(event Event) string {
	plain := renderPlain(event)
	// PR 1 preserves the current text path. PR 2 moves colors into this function.
	return plain
}

func renderTaskHeaderPlain(event Event) string {
	param := field(event, "param")
	if param != "" {
		param = fmt.Sprintf(" (%s)", param)
	}
	reason := field(event, "reason")
	if reason != "" {
		reason = fmt.Sprintf(" (%s)", reason)
	}
	return fmt.Sprintf("◼︎ %s%s%s\n", event.Text, param, reason)
}

func field(event Event, name string) string {
	for _, f := range event.Fields {
		if f.Name == name {
			return f.Value
		}
	}
	return ""
}

func joinArgs(args []string) string {
	return strings.Join(args, " ")
}

func cyan(v string) string { return color.Cyan(v).String() }
```

The `cyan` helper is intentionally unused in PR 1 and should be removed if lint complains; PR 2 will add actual color rendering.

- [ ] **Step 3: Add UI facade and testing constructor**

Create `pkg/ui/ui.go`:

```go
package ui

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type UI struct {
	out          io.Writer
	err          io.Writer
	renderer     Renderer
	debugEnabled bool
	events       []Event
}

func New(debugEnabled bool) *UI {
	return &UI{
		out:          os.Stdout,
		err:          os.Stderr,
		renderer:     TerminalRenderer{},
		debugEnabled: debugEnabled,
	}
}

func NewTesting(debugEnabled bool) (*bytes.Buffer, *UI) {
	buffer := bytes.NewBufferString("")
	return buffer, &UI{
		out:          buffer,
		err:          buffer,
		renderer:     PlainRenderer{},
		debugEnabled: debugEnabled,
	}
}

func (u *UI) Events() []Event {
	return append([]Event(nil), u.events...)
}

func (u *UI) SetOutputToStderr() {
	u.out = u.err
}

func (u *UI) emit(event Event) {
	u.events = append(u.events, event)
	Fprintf(u.out, "%s", u.renderer.Render(event))
}

func Fprintf(w io.Writer, format string, a ...any) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		log.Fatalf("failed to write to console: %s", err)
	}
}

func (u *UI) Debug(format string, params ...any) {
	if !u.debugEnabled {
		return
	}
	msg := strings.TrimSuffix(fmt.Sprintf(format, params...), "\n")
	u.emit(Event{Kind: KindDebug, Text: msg})
}

func (u *UI) Warningf(format string, params ...any) {
	u.emit(Event{Kind: KindWarning, Text: fmt.Sprintf(format, params...)})
}
```

- [ ] **Step 4: Move existing termui methods to emit events**

Modify `pkg/termui/ui.go`, `pkg/termui/task.go`, and `pkg/termui/hook.go` so `termui.UI` embeds or wraps `*ui.UI`. Keep the public method names unchanged.

The compatibility type should look like:

```go
package termui

import (
	"bytes"

	"github.com/devbuddy/devbuddy/pkg/config"
	baseui "github.com/devbuddy/devbuddy/pkg/ui"
)

type UI struct {
	*baseui.UI
}

func New(cfg *config.Config) *UI {
	return &UI{UI: baseui.New(cfg.DebugEnabled)}
}

func NewTesting(debugEnabled bool) (*bytes.Buffer, *UI) {
	buffer, ui := baseui.NewTesting(debugEnabled)
	return buffer, &UI{UI: ui}
}
```

Then implement each existing method by constructing the right event:

```go
func (u *UI) JumpProject(name string) {
	u.Emit(baseui.Event{Kind: baseui.KindJumpProject, Text: name})
}
```

If `Emit` needs to be public, rename `emit` to `Emit` in `pkg/ui/ui.go`.

- [ ] **Step 5: Preserve current tests and add event assertions**

Add a new test in `pkg/ui/ui_test.go`:

```go
func TestUIRecordsEventsWhileRenderingPlainOutput(t *testing.T) {
	buf, ui := NewTesting(false)

	ui.Emit(Event{Kind: KindJumpProject, Text: "org/repo"})

	require.Equal(t, "🐼  jumping to org/repo\n", buf.String())
	require.Equal(t, []Event{{Kind: KindJumpProject, Text: "org/repo"}}, ui.Events())
}
```

Run: `go test -count=1 ./pkg/ui ./pkg/termui ./pkg/context ./pkg/autoenv ./pkg/tasks/taskengine`

Expected: PASS.

- [ ] **Step 6: Commit and open PR 1**

Run:

```bash
script/test
script/lint
git add pkg/ui pkg/termui pkg/context pkg/autoenv pkg/tasks/taskengine
git commit -m "Introduce UI event facade"
git push -u origin ui-runtime-core
gh pr create --base main --head ui-runtime-core --title "Introduce UI event facade" --body-file /tmp/ui-runtime-core-pr.md
gh pr checks --watch --interval 10
```

Expected: CI passes before starting PR 2.

---

## PR 2: ANSI Policy Centralization

**Branch:** `ui-runtime-ansi`, based on `ui-runtime-core`.

**Files:**
- Modify: `pkg/ui/renderer.go`
- Modify: `pkg/ui/ui.go`
- Modify: `pkg/cmd/root.go`
- Modify: `pkg/helpers/open/open.go`
- Modify: `tests/internal/harness/output.go`
- Modify: `tests/internal/harness/strip_ansi.go` only if still needed
- Add tests in `pkg/ui/renderer_test.go`

- [ ] **Step 1: Add color mode**

Add to `pkg/ui/renderer.go`:

```go
type ColorMode string

const (
	ColorAuto   ColorMode = "auto"
	ColorAlways ColorMode = "always"
	ColorNever  ColorMode = "never"
)
```

Add `ColorMode` to `UI` and make `NewTesting` use `ColorNever`.

- [ ] **Step 2: Move aurora calls into TerminalRenderer**

Implement color in `renderTerminal` for existing event kinds. Keep the visible text equivalent to current output:

```go
func renderTerminal(event Event) string {
	switch event.Kind {
	case KindJumpProject:
		return fmt.Sprintf("🐼  %s %s\n", color.Yellow("jumping to"), color.Green(event.Text))
	case KindWarning:
		return fmt.Sprintf("%s: %s\n", color.Bold(color.Yellow("WARNING")), event.Text)
	}
	return renderPlain(event)
}
```

Move all remaining `aurora` imports out of `pkg/termui`.

- [ ] **Step 3: Route root errors through ui renderer or a small error renderer**

Modify `pkg/cmd/root.go` so the red `Error:` formatting uses the same color policy. If creating a root UI is too invasive, add `ui.RenderError(err)` and call it from `Execute`.

- [ ] **Step 4: Reduce ANSI stripping in tests**

Update tests that use `termui.NewTesting` so they assert plain output directly. Keep `StripAnsi` only for shell/PTY integration output, where command output may still include external ANSI.

- [ ] **Step 5: Verify**

Run:

```bash
script/test
script/lint
```

Expected: PASS. Shell tests are not required for this PR unless root error output changes.

- [ ] **Step 6: Commit and open PR 2**

Run:

```bash
git add pkg/ui pkg/termui pkg/cmd/root.go pkg/helpers/open/open.go tests/internal/harness
git commit -m "Centralize UI color rendering"
git push -u origin ui-runtime-ansi
gh pr create --base ui-runtime-core --head ui-runtime-ansi --title "Centralize UI color rendering" --body-file /tmp/ui-runtime-ansi-pr.md
gh pr checks --watch --interval 10
```

Expected: CI passes before starting PR 3.

---

## PR 3: Prompt Interface And Unit-Testable Worktree Switch

**Branch:** `ui-runtime-prompts`, based on `ui-runtime-ansi`.

**Files:**
- Create: `pkg/ui/prompts.go`
- Create: `pkg/ui/fake_prompts.go`
- Modify: `pkg/ui/ui.go`
- Modify: `pkg/cmd/worktree.go`
- Modify/Add: `pkg/cmd/worktree_test.go`
- Keep: `pkg/cmd/init.go` and `pkg/cmd/clone.go` on `promptui` for now

- [ ] **Step 1: Define prompt request types**

Create `pkg/ui/prompts.go`:

```go
package ui

import "errors"

var ErrPromptCancelled = errors.New("prompt cancelled")

type SelectOption struct {
	Value string
	Label string
}

type SelectRequest struct {
	Label   string
	Options []SelectOption
}

type ConfirmRequest struct {
	Label string
}

type Prompts interface {
	Select(SelectRequest) (string, error)
	Confirm(ConfirmRequest) (bool, error)
}
```

- [ ] **Step 2: Add fake prompts**

Create `pkg/ui/fake_prompts.go`:

```go
package ui

type FakePrompts struct {
	SelectRequests  []SelectRequest
	ConfirmRequests []ConfirmRequest
	SelectValue     string
	SelectErr       error
	ConfirmValue    bool
	ConfirmErr      error
}

func (p *FakePrompts) Select(req SelectRequest) (string, error) {
	p.SelectRequests = append(p.SelectRequests, req)
	return p.SelectValue, p.SelectErr
}

func (p *FakePrompts) Confirm(req ConfirmRequest) (bool, error) {
	p.ConfirmRequests = append(p.ConfirmRequests, req)
	return p.ConfirmValue, p.ConfirmErr
}
```

- [ ] **Step 3: Attach prompts to UI**

Add a `prompts Prompts` field to `pkg/ui.UI` plus:

```go
func (u *UI) Prompts() Prompts {
	return u.prompts
}

func (u *UI) SetPrompts(prompts Prompts) {
	u.prompts = prompts
}
```

`NewTesting` should install a `FakePrompts`.

- [ ] **Step 4: Refactor worktree switch selection to accept prompts**

In `pkg/cmd/worktree.go`, split selection into a pure function:

```go
func selectWorktree(prompts ui.Prompts, exec *executor.Executor, worktrees []worktree.Worktree) (worktree.Worktree, error)
```

Build `ui.SelectRequest` from existing `formatWorktreeRows` labels. Use the worktree path as `SelectOption.Value` so selection is stable even if labels duplicate.

- [ ] **Step 5: Add unit tests for worktree switch prompt behavior**

In `pkg/cmd/worktree_test.go`, add a test that calls the selection helper with a fake prompt and asserts:

```go
require.Equal(t, "Select worktree", fake.SelectRequests[0].Label)
require.Equal(t, selectedPath, got.Path)
```

Add a cancellation test:

```go
fake.SelectErr = ui.ErrPromptCancelled
_, err := selectWorktree(fake, exec, worktrees)
require.ErrorIs(t, err, ui.ErrPromptCancelled)
```

- [ ] **Step 6: Verify**

Run:

```bash
go test -count=1 ./pkg/cmd ./pkg/ui ./pkg/worktree
script/lint
```

Expected: PASS.

- [ ] **Step 7: Commit and open PR 3**

Run:

```bash
git add pkg/ui pkg/cmd/worktree.go pkg/cmd/worktree_test.go
git commit -m "Add prompt abstraction for worktree switching"
git push -u origin ui-runtime-prompts
gh pr create --base ui-runtime-ansi --head ui-runtime-prompts --title "Add prompt abstraction for worktree switching" --body-file /tmp/ui-runtime-prompts-pr.md
gh pr checks --watch --interval 10
```

Expected: CI passes before starting PR 4.

---

## PR 4: huh Prompt Adapter For Worktree Switch

**Branch:** `ui-runtime-huh-worktree`, based on `ui-runtime-prompts`.

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`
- Create: `pkg/ui/huh_prompts.go`
- Modify: `pkg/ui/ui.go`
- Modify: `pkg/cmd/worktree.go`
- Modify: `tests/shell/cmd_worktree_test.go`

- [ ] **Step 1: Add `huh` dependency**

Run:

```bash
go get github.com/charmbracelet/huh@v1.0.0
```

Expected: `go.mod` and `go.sum` update.

- [ ] **Step 2: Implement production prompts**

Create `pkg/ui/huh_prompts.go`:

```go
package ui

import (
	"errors"

	"github.com/charmbracelet/huh"
)

type HuhPrompts struct{}

func (HuhPrompts) Select(req SelectRequest) (string, error) {
	var selected string
	options := make([]huh.Option[string], 0, len(req.Options))
	for _, opt := range req.Options {
		options = append(options, huh.NewOption(opt.Label, opt.Value))
	}
	err := huh.NewSelect[string]().
		Title(req.Label).
		Options(options...).
		Value(&selected).
		Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return "", ErrPromptCancelled
	}
	return selected, err
}

func (HuhPrompts) Confirm(req ConfirmRequest) (bool, error) {
	var confirmed bool
	err := huh.NewConfirm().
		Title(req.Label).
		Value(&confirmed).
		Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return false, ErrPromptCancelled
	}
	return confirmed, err
}
```

- [ ] **Step 3: Use HuhPrompts in production**

Modify `pkg/ui.New`:

```go
prompts: HuhPrompts{},
```

Keep `NewTesting` on `FakePrompts`.

- [ ] **Step 4: Update worktree switch cancellation**

In `worktreeSwitchRun`, treat `ui.ErrPromptCancelled` as a no-op:

```go
if errors.Is(err, ui.ErrPromptCancelled) {
	return nil
}
```

- [ ] **Step 5: Add one PTY smoke test for ESC**

Add a shell test only in this PR:

```go
func Test_Cmd_Tree_Switch_Escape_Cancels(t *testing.T) {
	c := harness.NewDockerPTYInit(t)
	projectPath := "/home/tester/src/github.com/orgname/projname"
	worktreePath := "/home/tester/src/github.com/orgname/projname--feature-a"

	createGitProject(t, c, projectPath)
	c.Cd(t, projectPath)
	c.Run(t, "bud tree new feature-a")
	c.Cd(t, projectPath)

	c.Send(t, "bud tree switch\n")
	c.Expect(t, "Select worktree")
	c.Send(t, "\x1b")
	c.WaitPrompt(t)

	require.Equal(t, projectPath, c.Cwd(t))
	c.AssertExist(t, worktreePath+"/dev.yml")
}
```

If `huh` renders a different title string, adjust `c.Expect` to the actual stable visible text.

- [ ] **Step 6: Verify**

Run:

```bash
go test -count=1 ./pkg/ui ./pkg/cmd
TEST_DOCKER_IMAGE="ghcr.io/devbuddy/docker-testing:sha-7fd13f4" TEST_SHELL=bash go test -run Test_Cmd_Tree_Switch -count=1 ./tests/shell
TEST_DOCKER_IMAGE="ghcr.io/devbuddy/docker-testing:sha-7fd13f4" TEST_SHELL=zsh go test -run Test_Cmd_Tree_Switch -count=1 ./tests/shell
script/lint
```

Expected: PASS.

- [ ] **Step 7: Commit and open PR 4**

Run:

```bash
git add go.mod go.sum pkg/ui pkg/cmd/worktree.go tests/shell/cmd_worktree_test.go
git commit -m "Use huh for worktree switch prompt"
git push -u origin ui-runtime-huh-worktree
gh pr create --base ui-runtime-prompts --head ui-runtime-huh-worktree --title "Use huh for worktree switch prompt" --body-file /tmp/ui-runtime-huh-worktree-pr.md
gh pr checks --watch --interval 10
```

Expected: CI passes before starting PR 5.

---

## PR 5: Migrate Remaining Prompts And Remove promptui

**Branch:** `ui-runtime-remove-promptui`, based on `ui-runtime-huh-worktree`.

**Files:**
- Modify: `pkg/cmd/init.go`
- Modify: `pkg/cmd/clone.go`
- Modify: `pkg/cmd/worktree.go` prune confirmation
- Modify: `go.mod`
- Modify: `go.sum`
- Modify/add tests for init/clone/prune prompt behavior

- [ ] **Step 1: Migrate `bud init` template selection**

Replace direct `promptui.Select` usage with `ctx.UI.Prompts().Select` or a command-local `ui.Prompts` dependency. Convert selected template labels to stable template names.

- [ ] **Step 2: Migrate clone confirmation**

Replace `promptui.Prompt{IsConfirm: true}` with `ConfirmRequest`. Treat `ui.ErrPromptCancelled` and `false` as "do not create manifest".

- [ ] **Step 3: Migrate worktree prune confirmation**

Replace `promptui.Prompt{IsConfirm: true}` with `ConfirmRequest`. Keep existing behavior: skip deletion on cancellation/negative answer, continue to the next worktree.

- [ ] **Step 4: Remove promptui dependency**

Run:

```bash
go mod tidy
rg -n "promptui|manifoldco" .
```

Expected: no production references remain. `go.mod` no longer requires `github.com/manifoldco/promptui`.

- [ ] **Step 5: Verify**

Run:

```bash
script/test
script/lint
TEST_DOCKER_IMAGE="ghcr.io/devbuddy/docker-testing:sha-7fd13f4" TEST_SHELL=bash go test -run 'Test_Cmd_(Init|Clone|Tree_Prune)' -count=1 ./tests/shell
TEST_DOCKER_IMAGE="ghcr.io/devbuddy/docker-testing:sha-7fd13f4" TEST_SHELL=zsh go test -run 'Test_Cmd_(Init|Clone|Tree_Prune)' -count=1 ./tests/shell
```

Expected: PASS.

- [ ] **Step 6: Commit and open PR 5**

Run:

```bash
git add go.mod go.sum pkg/cmd pkg/ui tests
git commit -m "Migrate prompts to UI runtime"
git push -u origin ui-runtime-remove-promptui
gh pr create --base ui-runtime-huh-worktree --head ui-runtime-remove-promptui --title "Migrate prompts to UI runtime" --body-file /tmp/ui-runtime-remove-promptui-pr.md
gh pr checks --watch --interval 10
```

Expected: CI passes.

---

## Final Cleanup After Stack Merges

- Remove `pkg/termui` only after all call sites use `pkg/ui` directly.
- Revisit `tests/internal/harness/StripAnsi`; keep it only for PTY/shell integration tests that may include external command ANSI output.
- Update `ROADMAP.md` to mark terminal I/O extraction as underway or complete after the final PR lands.
- Consider adding `BUD_COLOR=auto|always|never` only after the renderer split is stable.
