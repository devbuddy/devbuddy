package termui

import (
	"fmt"
	"strings"

	color "github.com/logrusorgru/aurora"
)

func (u *UI) TaskHeader(name, param, reason string) {
	if param != "" {
		param = fmt.Sprintf(" (%s)", color.Blue(param))
	}
	if reason != "" {
		reason = fmt.Sprintf(" (%s)", color.Brown(reason))
	}
	u.printf("%s %s%s%s\n", color.Brown("◼︎"), color.Magenta(name), param, reason)
}

func (u *UI) TaskCommand(cmdline string, args ...string) {
	u.printf("  Running: %s %s\n", color.Bold(color.Cyan(cmdline)), color.Cyan(strings.Join(args, " ")))
}

func (u *UI) TaskShell(cmdline string) {
	u.printf("  Running: %s\n", color.Cyan(cmdline))
}

func (u *UI) TaskActed() {
	u.printf("  %s\n", color.Green("Done!"))
}

func (u *UI) TaskAlreadyOk() {
	u.printf("  %s\n", color.Green("Already OK!"))
}

func (u *UI) TaskError(err error) {
	u.printf("  %s\n", color.Red(err.Error()))
}

func (u *UI) TaskErrorf(message string, a ...interface{}) {
	u.TaskError(fmt.Errorf(message, a...))
}

func (u *UI) TaskWarning(message string) {
	u.printf("  Warning: %s\n", color.Brown(message))
}

func (u *UI) TaskActionHeader(desc string) {
	u.printf("  %s%s\n", color.Brown("▪︎"), color.Magenta(desc))
}

func (u *UI) ActionHeader(description string) {
	u.printf("🐼  %s\n", color.Cyan(description))
}

func (u *UI) ActionNotice(text string) {
	u.printf("⚠️   %s\n", color.Brown(text))
}

func (u *UI) ActionDone() {
	u.printf("✅  %s\n", color.Green("Done!"))
}
