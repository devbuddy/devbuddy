package termui

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/config"
	baseui "github.com/devbuddy/devbuddy/pkg/ui"
)

func Fprintf(w io.Writer, format string, a ...any) {
	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		log.Fatalf("failed to write to console: %s", err)
	}
}

type UI struct {
	out          io.Writer
	err          io.Writer
	debugEnabled bool
	recorder     *baseui.UI
	renderer     baseui.Renderer
}

func New(cfg *config.Config) *UI {
	return &UI{
		out:          os.Stdout,
		err:          os.Stderr,
		debugEnabled: cfg.DebugEnabled,
		recorder:     baseui.New(),
		renderer:     baseui.TerminalRenderer{},
	}
}

func NewTesting(debugEnabled bool) (*bytes.Buffer, *UI) {
	buffer := bytes.NewBufferString("")
	return buffer, &UI{
		out:          buffer,
		err:          buffer,
		debugEnabled: debugEnabled,
		recorder:     baseui.New(),
		renderer:     baseui.PlainRenderer{},
	}
}

func (u *UI) Events() []baseui.Event {
	return u.recorder.Events()
}

func (u *UI) record(event baseui.Event) {
	u.recorder.Record(event)
}

func (u *UI) emit(event baseui.Event) {
	u.record(event)
	Fprintf(u.out, "%s", u.renderer.Render(event))
}

func (u *UI) emitErr(event baseui.Event) {
	u.record(event)
	Fprintf(u.err, "%s", u.renderer.Render(event))
}

func (u *UI) SetOutputToStderr() {
	u.out = u.err
}

func (u *UI) Debug(format string, params ...any) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		msg = strings.TrimSuffix(msg, "\n")
		u.emit(baseui.Event{Kind: baseui.KindDebug, Text: msg})
	}
}

func (u *UI) Warningf(format string, params ...any) {
	msg := fmt.Sprintf(format, params...)
	u.emit(baseui.Event{Kind: baseui.KindWarning, Text: msg})
}

func (u *UI) CommandHeader(cmdline string) {
	u.emitErr(baseui.Event{Kind: baseui.KindCommandHeader, Text: cmdline})
}

func (u *UI) CommandRun(cmdline string, args ...string) {
	u.emit(baseui.Event{Kind: baseui.KindCommandRun, Text: cmdline, Fields: []baseui.Field{baseui.F("args", strings.Join(args, " "))}})
}

func (u *UI) CommandActed() {
	u.emit(baseui.Event{Kind: baseui.KindCommandActed})
}

func (u *UI) ProjectExists() {
	u.emit(baseui.Event{Kind: baseui.KindProjectExists})
}

func (u *UI) JumpProject(name string) {
	u.emit(baseui.Event{Kind: baseui.KindJumpProject, Text: name})
}
