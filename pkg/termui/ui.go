package termui

import (
	"fmt"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/config"
)

type UI struct {
	baseUI
}

func NewUI(cfg *config.Config) *UI {
	return &UI{
		newBaseUI(cfg),
	}
}

func (u *UI) TaskHeader(name string, param string) {
	fmt.Fprintf(u.out, "%s %s (%s)\n", color.Brown("★"), color.Magenta(name), color.Gray(param))
}

func (u *UI) TaskActed() {
	fmt.Fprintf(u.out, "  %s\n", color.Green("Done!"))
}

func (u *UI) TaskAlreadyOk() {
	fmt.Fprintf(u.out, "  %s\n", color.Green("Already OK!"))
}

func (u *UI) TaskError(err error) {
	fmt.Fprintf(u.out, "  %s %s\n", color.Red("Error:"), color.Brown(err.Error()))
}
