package termui

import (
	"fmt"
	"os"
	"strings"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/config"
)

type UI struct {
	baseUI
}

func NewUI(cfg *config.Config) *UI {
	return &UI{
		baseUI{
			out:          os.Stdout,
			debugEnabled: cfg.DebugEnabled,
		},
	}
}

func (u *UI) CommandHeader(cmdline string) {
	fmt.Fprintf(os.Stderr, "🐼  %s %s\n", color.Blue("running"), color.Cyan(cmdline))
}

func (u *UI) TaskHeader(name string, param string) {
	if param != "" {
		param = fmt.Sprintf(" (%s)", color.Blue(param))
	}
	fmt.Fprintf(u.out, "%s %s%s\n", color.Brown("◼︎"), color.Magenta(name), param)
}

func (u *UI) TaskCommand(cmdline string, args ...string) {
	fmt.Fprint(u.out, "  Running: ", color.Bold(color.Cyan(cmdline)), color.Cyan(strings.Join(args, " ")))
}

func (u *UI) TaskActed() {
	fmt.Fprintf(u.out, "  %s\n", color.Green("Done!"))
}

func (u *UI) TaskAlreadyOk() {
	fmt.Fprintf(u.out, "  %s\n", color.Green("Already OK!"))
}

func (u *UI) TaskError(err error) {
	fmt.Fprintf(u.out, "  %s\n", color.Red(err.Error()))
}

func (u *UI) TaskWarning(message string) {
	fmt.Fprintf(u.out, "  Warning: %s\n", color.Brown(message))
}

func (u *UI) ProjectExists() {
	fmt.Fprintf(u.out, "🐼  %s\n", color.Brown("project already exists locally"))
}

func (u *UI) JumpProject(name string) {
	fmt.Fprintf(u.out, "🐼  %s %s\n", color.Brown("jumping to"), color.Green(name))
}
