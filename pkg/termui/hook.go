package termui

import (
	"fmt"

	color "github.com/logrusorgru/aurora"
)

type HookUI struct {
	baseUI
}

func NewHookUI() *HookUI {
	return &HookUI{
		newBaseUI(),
	}
}

func (u *HookUI) HookFeatureActivated(name string, version string) {
	msg := fmt.Sprintf("%s(%s) activated.", name, version)
	fmt.Fprintf(u.out, "🐼  %s\n", color.Cyan(msg))
}

func (u *HookUI) HookFeatureFailure(name string, version string) {
	msg := fmt.Sprintf("failed to activate %s(%s). Try running dad up first!", name, version)
	fmt.Fprintf(u.out, "🐼  %s\n", color.Red(msg))
}
