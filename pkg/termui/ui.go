package termui

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	color "github.com/logrusorgru/aurora"

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
	debugEnabled bool
	recorder     *baseui.UI
}

func New(cfg *config.Config) *UI {
	return &UI{
		out:          os.Stdout,
		debugEnabled: cfg.DebugEnabled,
		recorder:     baseui.New(),
	}
}

func NewTesting(debugEnabled bool) (*bytes.Buffer, *UI) {
	buffer := bytes.NewBufferString("")
	return buffer, &UI{
		out:          buffer,
		debugEnabled: debugEnabled,
		recorder:     baseui.New(),
	}
}

func (u *UI) Events() []baseui.Event {
	return u.recorder.Events()
}

func (u *UI) record(event baseui.Event) {
	u.recorder.Record(event)
}

func (u *UI) SetOutputToStderr() {
	u.out = os.Stderr
}

func (u *UI) Debug(format string, params ...any) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		msg = strings.TrimSuffix(msg, "\n")
		u.record(baseui.Event{Kind: baseui.KindDebug, Text: msg})
		Fprintf(u.out, "%s: %s\n", color.Yellow("BUD_DEBUG"), color.Gray(12, msg))
	}
}

func (u *UI) Warningf(format string, params ...any) {
	msg := fmt.Sprintf(format, params...)
	u.record(baseui.Event{Kind: baseui.KindWarning, Text: msg})
	Fprintf(u.out, "%s: %s\n", color.Bold(color.Yellow("WARNING")), msg)
}

func (u *UI) CommandHeader(cmdline string) {
	u.record(baseui.Event{Kind: baseui.KindCommandHeader, Text: cmdline})
	Fprintf(os.Stderr, "🐼  %s %s\n", color.Blue("running"), color.Cyan(cmdline))
}

func (u *UI) CommandRun(cmdline string, args ...string) {
	u.record(baseui.Event{Kind: baseui.KindCommandRun, Text: cmdline, Fields: []baseui.Field{baseui.F("args", strings.Join(args, " "))}})
	Fprintf(u.out, "%s %s\n", color.Bold(color.Cyan(cmdline)), color.Cyan(strings.Join(args, " ")))
}

func (u *UI) CommandActed() {
	u.record(baseui.Event{Kind: baseui.KindCommandActed})
	Fprintf(u.out, "  %s\n", color.Green("Done!"))
}

func (u *UI) ProjectExists() {
	u.record(baseui.Event{Kind: baseui.KindProjectExists})
	Fprintf(u.out, "🐼  %s\n", color.Yellow("project already exists locally"))
}

func (u *UI) JumpProject(name string) {
	u.record(baseui.Event{Kind: baseui.KindJumpProject, Text: name})
	Fprintf(u.out, "🐼  %s %s\n", color.Yellow("jumping to"), color.Green(name))
}
