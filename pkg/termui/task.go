package termui

import (
	"fmt"
	"strings"

	baseui "github.com/devbuddy/devbuddy/pkg/ui"
)

func (u *UI) TaskHeader(name, param, reason string) {
	u.emit(baseui.Event{Kind: baseui.KindTaskHeader, Text: name, Fields: []baseui.Field{baseui.F("param", param), baseui.F("reason", reason)}})
}

func (u *UI) TaskCommand(cmdline string, args ...string) {
	u.emit(baseui.Event{Kind: baseui.KindTaskCommand, Text: cmdline, Fields: []baseui.Field{baseui.F("args", strings.Join(args, " "))}})
}

func (u *UI) TaskShell(cmdline string) {
	u.emit(baseui.Event{Kind: baseui.KindTaskShell, Text: cmdline})
}

func (u *UI) TaskActed() {
	u.emit(baseui.Event{Kind: baseui.KindTaskActed})
}

func (u *UI) TaskAlreadyOk() {
	u.emit(baseui.Event{Kind: baseui.KindTaskAlreadyOK})
}

func (u *UI) TaskError(err error) {
	u.emit(baseui.Event{Kind: baseui.KindTaskError, Text: err.Error()})
}

func (u *UI) TaskErrorf(message string, a ...interface{}) {
	u.TaskError(fmt.Errorf(message, a...))
}

func (u *UI) TaskWarning(message string) {
	u.emit(baseui.Event{Kind: baseui.KindTaskWarning, Text: message})
}

func (u *UI) TaskActionHeader(desc string) {
	u.emit(baseui.Event{Kind: baseui.KindTaskActionHeader, Text: desc})
}

func (u *UI) ActionHeader(description string) {
	u.emit(baseui.Event{Kind: baseui.KindActionHeader, Text: description})
}

func (u *UI) ActionNotice(text string) {
	u.emit(baseui.Event{Kind: baseui.KindActionNotice, Text: text})
}

func (u *UI) ActionDone() {
	u.emit(baseui.Event{Kind: baseui.KindActionDone})
}
