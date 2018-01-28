package termui

import (
	"fmt"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/config"
)

type HookUI struct {
	baseUI
}

func NewHookUI(cfg *config.Config) *HookUI {
	return &HookUI{
		newBaseUI(cfg),
	}
}

func (u *HookUI) HookFeatureActivated(name string, version string) {
	msg := color.Sprintf("%s activated.", name)
	ver := color.Sprintf("(version: %s)", version)
	fmt.Fprintf(u.out, "🐼  %s %s\n", color.Cyan(msg), color.Blue(ver))
}

func (u *HookUI) HookFeatureFailure(name string, version string) {
	msg := color.Sprintf("failed to activate %s. Try running dad up first!", name)
	ver := color.Sprintf("(version: %s)", version)
	fmt.Fprintf(u.out, "🐼  %s %s\n", color.Red(msg), color.Brown(ver))
}
