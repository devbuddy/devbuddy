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
)

func Fprintf(w io.Writer, format string, a ...interface{}) {
	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		log.Fatalf("failed to write to console: %s", err)
	}
}

type UI struct {
	out          io.Writer
	debugEnabled bool
}

func New(cfg *config.Config) *UI {
	return &UI{
		out:          os.Stdout,
		debugEnabled: cfg.DebugEnabled,
	}
}

func NewHook(cfg *config.Config) *UI {
	return &UI{
		out:          os.Stderr,
		debugEnabled: cfg.DebugEnabled,
	}
}

func NewTesting(debugEnabled bool) (*bytes.Buffer, *UI) {
	buffer := bytes.NewBufferString("")
	return buffer, &UI{
		out:          buffer,
		debugEnabled: debugEnabled,
	}
}

func (u *UI) Debug(format string, params ...interface{}) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		Fprintf(u.out, "BUD_DEBUG: %s\n", color.Gray(msg))
	}
}

func (u *UI) CommandHeader(cmdline string) {
	Fprintf(os.Stderr, "🐼  %s %s\n", color.Blue("running"), color.Cyan(cmdline))
}

func (u *UI) CommandRun(cmdline string, args ...string) {
	Fprintf(u.out, "%s %s\n", color.Bold(color.Cyan(cmdline)), color.Cyan(strings.Join(args, " ")))
}

func (u *UI) CommandActed() {
	Fprintf(u.out, "  %s\n", color.Green("Done!"))
}

func (u *UI) ProjectExists() {
	Fprintf(u.out, "🐼  %s\n", color.Brown("project already exists locally"))
}

func (u *UI) JumpProject(name string) {
	Fprintf(u.out, "🐼  %s %s\n", color.Brown("jumping to"), color.Green(name))
}
