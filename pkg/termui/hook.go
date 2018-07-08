package termui

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"

	"github.com/devbuddy/devbuddy/pkg/config"
)

type HookUI struct {
	baseUI
}

func NewHookUI(cfg *config.Config) *HookUI {
	return &HookUI{
		baseUI{
			out:          os.Stderr,
			debugEnabled: cfg.DebugEnabled,
		},
	}
}

func (u *HookUI) HookFeatureActivated(name string, version string) {
	msg := color.Sprintf("%s activated.", name)
	ver := color.Sprintf("(version: %s)", version)
	fmt.Fprintf(u.out, "🐼  %s %s\n", color.Cyan(msg), color.Blue(ver))
}

func (u *HookUI) HookFeatureFailure(name string, version string) {
	msg := color.Sprintf("failed to activate %s. Try running 'bud up' first!", name)
	ver := color.Sprintf("(version: %s)", version)
	fmt.Fprintf(u.out, "🐼  %s %s\n", color.Red(msg), color.Brown(ver))
}

func HookShellDetectionError(err error) {
	fmt.Fprintln(os.Stderr, color.Brown("Could not detect your shell:"), err.Error())
}

func HookIntegrationError(message string) {
	fmt.Fprintf(os.Stderr, "%s: %s", color.Red("Shell integration error:"), message)
}
