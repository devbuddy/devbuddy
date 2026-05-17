package ui

import (
	"fmt"
	"strings"
)

func (u *UI) TaskHeader(name, param, reason string) {
	u.emit(Event{Kind: KindTaskHeader, Text: name, Fields: []Field{F("param", param), F("reason", reason)}})
}

func (u *UI) TaskCommand(cmdline string, args ...string) {
	u.emit(Event{Kind: KindTaskCommand, Text: cmdline, Fields: []Field{F("args", strings.Join(args, " "))}})
}

func (u *UI) TaskShell(cmdline string) {
	u.emit(Event{Kind: KindTaskShell, Text: cmdline})
}

func (u *UI) TaskActed() {
	u.emit(Event{Kind: KindTaskActed})
}

func (u *UI) TaskAlreadyOk() {
	u.emit(Event{Kind: KindTaskAlreadyOK})
}

func (u *UI) TaskError(err error) {
	u.emit(Event{Kind: KindTaskError, Text: err.Error()})
}

func (u *UI) TaskErrorf(message string, a ...interface{}) {
	u.TaskError(fmt.Errorf(message, a...))
}

func (u *UI) TaskWarning(message string) {
	u.emit(Event{Kind: KindTaskWarning, Text: message})
}

func (u *UI) TaskActionHeader(desc string) {
	u.emit(Event{Kind: KindTaskActionHeader, Text: desc})
}

func (u *UI) ActionHeader(description string) {
	u.emit(Event{Kind: KindActionHeader, Text: description})
}

func (u *UI) ActionNotice(text string) {
	u.emit(Event{Kind: KindActionNotice, Text: text})
}

func (u *UI) ActionDone() {
	u.emit(Event{Kind: KindActionDone})
}
