package termui

import (
	"fmt"
	"strings"

	color "github.com/logrusorgru/aurora"

	baseui "github.com/devbuddy/devbuddy/pkg/ui"
)

func (u *UI) TaskHeader(name, param, reason string) {
	u.record(baseui.Event{Kind: baseui.KindTaskHeader, Text: name, Fields: []baseui.Field{baseui.F("param", param), baseui.F("reason", reason)}})
	if param != "" {
		param = fmt.Sprintf(" (%s)", color.Blue(param))
	}
	if reason != "" {
		reason = fmt.Sprintf(" (%s)", color.Yellow(reason))
	}
	Fprintf(u.out, "%s %s%s%s\n", color.Yellow("◼︎"), color.Magenta(name), param, reason)
}

func (u *UI) TaskCommand(cmdline string, args ...string) {
	u.record(baseui.Event{Kind: baseui.KindTaskCommand, Text: cmdline, Fields: []baseui.Field{baseui.F("args", strings.Join(args, " "))}})
	Fprintf(u.out, "  Running: %s %s\n", color.Bold(color.Cyan(cmdline)), color.Cyan(strings.Join(args, " ")))
}

func (u *UI) TaskShell(cmdline string) {
	u.record(baseui.Event{Kind: baseui.KindTaskShell, Text: cmdline})
	Fprintf(u.out, "  Running: %s\n", color.Cyan(cmdline))
}

func (u *UI) TaskActed() {
	u.record(baseui.Event{Kind: baseui.KindTaskActed})
	Fprintf(u.out, "  %s\n", color.Green("Done!"))
}

func (u *UI) TaskAlreadyOk() {
	u.record(baseui.Event{Kind: baseui.KindTaskAlreadyOK})
	Fprintf(u.out, "  %s\n", color.Green("Already OK!"))
}

func (u *UI) TaskError(err error) {
	u.record(baseui.Event{Kind: baseui.KindTaskError, Text: err.Error()})
	Fprintf(u.out, "  %s\n", color.Red(err.Error()))
}

func (u *UI) TaskErrorf(message string, a ...interface{}) {
	u.TaskError(fmt.Errorf(message, a...))
}

func (u *UI) TaskWarning(message string) {
	u.record(baseui.Event{Kind: baseui.KindTaskWarning, Text: message})
	Fprintf(u.out, "  Warning: %s\n", color.Yellow(message))
}

func (u *UI) TaskActionHeader(desc string) {
	u.record(baseui.Event{Kind: baseui.KindTaskActionHeader, Text: desc})
	Fprintf(u.out, "  %s%s\n", color.Yellow("▪︎"), color.Magenta(desc))
}

func (u *UI) ActionHeader(description string) {
	u.record(baseui.Event{Kind: baseui.KindActionHeader, Text: description})
	Fprintf(u.out, "🐼  %s\n", color.Cyan(description))
}

func (u *UI) ActionNotice(text string) {
	u.record(baseui.Event{Kind: baseui.KindActionNotice, Text: text})
	Fprintf(u.out, "⚠️   %s\n", color.Yellow(text))
}

func (u *UI) ActionDone() {
	u.record(baseui.Event{Kind: baseui.KindActionDone})
	Fprintf(u.out, "✅  %s\n", color.Green("Done!"))
}
