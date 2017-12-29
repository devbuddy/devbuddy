package termui

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"
)

type UI struct {
	out *os.File
}

func NewHookUI() *UI {
	return &UI{
		out: os.Stderr,
	}
}

func (u *UI) HookFeature(format string, params ...interface{}) {
	msg := fmt.Sprintf(format, params...)
	fmt.Fprintf(u.out, "🐼  %s\n", color.Cyan(msg))
}

func (u *UI) Debug(format string, params ...interface{}) {
	msg := fmt.Sprintf(format, params...)
	fmt.Fprintf(u.out, "DEBUG: %s\n", color.Gray(msg))
}
