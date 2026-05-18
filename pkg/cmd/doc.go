package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"

	devbuddy "github.com/devbuddy/devbuddy"
)

func printBudDoc(w io.Writer) error {
	if !isTerminal(w) {
		_, err := fmt.Fprint(w, devbuddy.Documentation)
		return err
	}

	rendered, err := renderBudDoc(devbuddy.Documentation)
	if err != nil {
		rendered = devbuddy.Documentation
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
