package ui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type UI struct {
	out          io.Writer
	err          io.Writer
	debugEnabled bool
	renderer     Renderer
	events       []Event
	prompts      Prompts
}

func New() *UI {
	return &UI{
		out:      os.Stdout,
		err:      os.Stderr,
		renderer: TerminalRenderer{},
		prompts:  SurveyPrompts{},
	}
}

func NewTerminal(debugEnabled bool) *UI {
	ui := New()
	ui.debugEnabled = debugEnabled
	return ui
}

func NewTesting() (*FakePrompts, *UI) {
	_, ui := NewBufferedTesting(false)
	return ui.prompts.(*FakePrompts), ui
}

func NewBufferedTesting(debugEnabled bool) (*bytes.Buffer, *UI) {
	buffer := bytes.NewBufferString("")
	prompts := &FakePrompts{}
	return buffer, &UI{
		out:          buffer,
		err:          buffer,
		debugEnabled: debugEnabled,
		renderer:     PlainRenderer{},
		prompts:      prompts,
	}
}

func (u *UI) Record(event Event) {
	u.events = append(u.events, event)
}

func (u *UI) Events() []Event {
	return append([]Event(nil), u.events...)
}

func (u *UI) Prompts() Prompts {
	return u.prompts
}

func (u *UI) SetPrompts(prompts Prompts) {
	u.prompts = prompts
}

func (u *UI) emit(event Event) {
	u.Record(event)
	Fprintf(u.out, "%s", u.renderer.Render(event))
}

func (u *UI) emitErr(event Event) {
	u.Record(event)
	Fprintf(u.err, "%s", u.renderer.Render(event))
}

func (u *UI) SetOutputToStderr() {
	u.out = u.err
}

func (u *UI) Debug(format string, params ...any) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		msg = strings.TrimSuffix(msg, "\n")
		u.emit(Event{Kind: KindDebug, Text: msg})
	}
}

func (u *UI) Warningf(format string, params ...any) {
	msg := fmt.Sprintf(format, params...)
	u.emit(Event{Kind: KindWarning, Text: msg})
}

func (u *UI) CommandHeader(cmdline string) {
	u.emitErr(Event{Kind: KindCommandHeader, Text: cmdline})
}

func (u *UI) CommandRun(cmdline string, args ...string) {
	u.emit(Event{Kind: KindCommandRun, Text: cmdline, Fields: []Field{F("args", strings.Join(args, " "))}})
}

func (u *UI) CommandActed() {
	u.emit(Event{Kind: KindCommandActed})
}

func (u *UI) ProjectExists() {
	u.emit(Event{Kind: KindProjectExists})
}

func (u *UI) JumpProject(name string) {
	u.emit(Event{Kind: KindJumpProject, Text: name})
}
