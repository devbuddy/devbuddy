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
	debugBuffer  bytes.Buffer
}

func New(cfg *config.Config) *UI {
	return &UI{
		out:          os.Stdout,
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

func (u *UI) SetOutputToStderr() {
	u.out = os.Stderr
}

func (u *UI) Debug(format string, params ...interface{}) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		msg = strings.TrimSuffix(msg, "\n")
		u.printf("%s: %s\n", color.Brown("BUD_DEBUG"), color.Gray(msg))
	}
}

func (u *UI) FlushDebugBuffer() []byte {
	b := u.debugBuffer.Bytes()
	u.debugBuffer.Reset()
	return b
}

func (u *UI) printf(format string, params ...interface{}) {
	s := fmt.Sprintf(format, params...)
	u.debugBuffer.WriteString(s)

	_, err := u.out.Write([]byte(s))
	if err != nil {
		log.Fatalf("failed to write to console: %s", err)
	}
}

func (u *UI) Warningf(format string, params ...interface{}) {
	msg := fmt.Sprintf(format, params...)
	u.printf("%s: %s\n", color.Bold(color.Brown("WARNING")), msg)
}

func (u *UI) CommandHeader(cmdline string) {
	Fprintf(os.Stderr, "🐼  %s %s\n", color.Blue("running"), color.Cyan(cmdline))
}

func (u *UI) CommandRun(cmdline string, args ...string) {
	u.printf("%s %s\n", color.Bold(color.Cyan(cmdline)), color.Cyan(strings.Join(args, " ")))
}

func (u *UI) CommandActed() {
	u.printf("  %s\n", color.Green("Done!"))
}

func (u *UI) ProjectExists() {
	u.printf("🐼  %s\n", color.Brown("project already exists locally"))
}

func (u *UI) JumpProject(name string) {
	u.printf("🐼  %s %s\n", color.Brown("jumping to"), color.Green(name))
}
