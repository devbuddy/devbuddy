package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

const budDoc = `# DevBuddy Project Guide

## Purpose

DevBuddy reads dev.yml at the project root and prepares the whole development
environment: language runtimes, virtualenvs, system packages, environment
variables, project commands, project links, and shell activation.

## First Steps

- Run bud up from the project root before building, testing, or editing runtime
  assumptions. It installs or verifies the tasks listed in dev.yml.
- If shell integration is unavailable, bud up still performs installs and checks,
  but environment changes do not persist in the parent shell. Inspect dev.yml and
  prefer bud project commands for work that needs project env vars.
- Install shell integration with eval "$(bud --shell-init)" so DevBuddy
  automatically activates projects when cd-ing between them.
- Run bud inspect to see the detected project and parsed up tasks.
- Run bud --help from the project root to see project-specific commands from dev.yml.
  Also read the commands: section directly when choosing validation commands.

## dev.yml Sections

- env: environment variables DevBuddy exports when the project is active.
- up: setup tasks run by bud up. Common tasks include go, python, ruby, node,
  apt, homebrew, pip, pipfile, python_develop, envfile, and custom.
- commands: project-local commands exposed as bud <name>. Prefer these over
  guessing test, lint, or build commands.
- open: named project URLs available through bud open <name>.

## Using bud Commands

- bud up: install or verify the project environment.
- bud inspect: print the project path and parsed setup tasks.
- bud --help: list built-in commands and project-specific commands from dev.yml.
- bud <command>: run a project command from dev.yml with the project's env vars.
- bud open --list: list project URLs.
- bud cd <project>: jump to a known local project through shell finalizers.
- bud tree <subcommand>: manage Git worktrees when available.

## Guidance for Automation

- Treat dev.yml as the source of truth for setup and project commands.
- Prefer project-specific bud commands over ad hoc package-manager commands when
  available.
- Run the narrowest relevant bud command or validation command after edits.
- Avoid release, deploy, destructive cleanup, or broad system package changes
  unless the user explicitly asks.
- If bud up opens an OS installer or asks for manual action, stop and report the
  instruction instead of trying to bypass it.

## More Detail

- Run bud --help for CLI usage.
- Read docs/Config.md or https://github.com/devbuddy/devbuddy/blob/master/docs/Config.md
  for the full dev.yml reference.
`

func printBudDoc(w io.Writer) error {
	if !isTerminal(w) {
		_, err := fmt.Fprint(w, budDoc)
		return err
	}

	rendered, err := renderBudDoc(budDoc)
	if err != nil {
		rendered = budDoc
	}
	if pageErr := pageBudDoc(w, rendered); pageErr == nil {
		return nil
	}
	_, err = fmt.Fprint(w, rendered)
	return err
}

func renderBudDoc(markdown string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return "", err
	}
	return renderer.Render(markdown)
}

func pageBudDoc(w io.Writer, content string) error {
	pager := strings.TrimSpace(os.Getenv("PAGER"))
	if pager == "" {
		pager = "less -R"
	}

	cmd := exec.Command("sh", "-c", pager)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "LESS=FRX")
	return cmd.Run()
}

func isTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(file.Fd()))
}
