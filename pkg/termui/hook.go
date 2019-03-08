package termui

import (
	"os"

	color "github.com/logrusorgru/aurora"
)

func (u *UI) HookFeatureActivated(name string, version string) {
	msg := color.Sprintf("%s activated.", name)
	ver := color.Sprintf("(version: %s)", version)
	Fprintf(u.out, "🐼  %s %s\n", color.Cyan(msg), color.Blue(ver))
}

func (u *UI) HookFeatureFailure(name string, version string) {
	msg := color.Sprintf("failed to activate %s. Try running 'bud up' first!", name)
	ver := color.Sprintf("(version: %s)", version)
	Fprintf(u.out, "🐼  %s %s\n", color.Red(msg), color.Brown(ver))
}

func HookShellDetectionError(err error) {
	Fprintf(os.Stderr, "%s %s\n", color.Brown("Could not detect your shell:"), err.Error())
}
