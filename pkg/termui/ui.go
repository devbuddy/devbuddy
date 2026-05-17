package termui

import (
	"bytes"
	"io"

	"github.com/devbuddy/devbuddy/pkg/config"
	baseui "github.com/devbuddy/devbuddy/pkg/ui"
)

type UI = baseui.UI

func Fprintf(w io.Writer, format string, a ...any) {
	baseui.Fprintf(w, format, a...)
}

func New(cfg *config.Config) *UI {
	return baseui.NewTerminal(cfg.DebugEnabled)
}

func NewTesting(debugEnabled bool) (*bytes.Buffer, *UI) {
	return baseui.NewBufferedTesting(debugEnabled)
}
