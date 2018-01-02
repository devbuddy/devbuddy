package termui

import (
	"fmt"

	color "github.com/logrusorgru/aurora"
)

type UI struct {
	baseUI
}

func NewUI() *UI {
	return &UI{
		newBaseUI(),
	}
}

func (u *UI) TaskHeader(name string, param string) {
	fmt.Fprintf(u.out, "%s %s (%s)\n", color.Brown("★"), color.Magenta(name), color.Gray(param))
}
