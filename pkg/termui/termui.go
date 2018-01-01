package termui

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/config"
)

type UI struct {
	out          *os.File
	debugEnabled bool
}

func NewHookUI() *UI {
	return &UI{
		out:          os.Stderr,
		debugEnabled: config.DebugEnabled(),
	}
}

func (u *UI) HookFeatureActivated(name string, version string) {
	msg := fmt.Sprintf("%s(%s) activated.", name, version)
	fmt.Fprintf(u.out, "🐼  %s\n", color.Cyan(msg))
}

func (u *UI) HookFeatureFailure(name string, version string) {
	msg := fmt.Sprintf("failed to activate %s(%s). Try running dad up first!", name, version)
	fmt.Fprintf(u.out, "🐼  %s\n", color.Red(msg))
}

func (u *UI) Debug(format string, params ...interface{}) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		fmt.Fprintf(u.out, "DEBUG: %s\n", color.Gray(msg))
	}
}
