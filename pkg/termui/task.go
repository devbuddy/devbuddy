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
		reason = fmt.Sprintf(" (%s)", color.Gray(reason))
	}
	Fprintf(u.out, "%s %s%s%s\n", color.Brown("◼︎"), color.Magenta(name), param, reason)
}

func (u *UI) TaskCommand(cmdline string, args ...string) {
	Fprintf(u.out, "  Running: %s %s\n", color.Bold(color.Cyan(cmdline)), color.Cyan(strings.Join(args, " ")))
}

func (u *UI) TaskShell(cmdline string) {
	Fprintf(u.out, "  Running: %s\n", color.Cyan(cmdline))
}

func (u *UI) TaskActed() {
	Fprintf(u.out, "  %s\n", color.Green("Done!"))
}

func (u *UI) TaskAlreadyOk() {
	Fprintf(u.out, "  %s\n", color.Green("Already OK!"))
}

func (u *UI) TaskError(err error) {
	Fprintf(u.out, "  %s\n", color.Red(err.Error()))
}

func (u *UI) TaskWarning(message string) {
	Fprintf(u.out, "  Warning: %s\n", color.Brown(message))
}

func (u *UI) TaskActionHeader(desc string) {
	Fprintf(u.out, "  %s%s\n", color.Brown("▪︎"), color.Magenta(desc))
}

func (u *UI) ActionHeader(description string) {
	Fprintf(u.out, "🐼  %s\n", color.Cyan(description))
}

func (u *UI) ActionNotice(text string) {
	Fprintf(u.out, "⚠️   %s\n", color.Brown(text))
}

func (u *UI) ActionDone() {
	Fprintf(u.out, "✅  %s\n", color.Green("Done!"))
}
