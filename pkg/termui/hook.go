package termui

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"
)

func (u *UI) HookFeatureActivated(name string, param string) {
	msg := fmt.Sprintf("%s activated.", name)

	paramStr := ""
	if param != "" {
		paramStr = fmt.Sprintf(" (%s)", param)
	}

	Fprintf(u.out, "🐼  %s%s\n", color.Cyan(msg), color.Blue(paramStr))
}

func (u *UI) HookFeatureFailure(name string, param string) {
	msg := fmt.Sprintf("failed to activate %s. Try running 'bud up' first!", name)

	paramStr := ""
	if param != "" {
		paramStr = fmt.Sprintf(" (%s)", param)
	}

	Fprintf(u.out, "🐼  %s%s\n", color.Red(msg), color.Brown(paramStr))
}

func HookShellDetectionError(err error) {
	Fprintf(os.Stderr, "%s %s\n", color.Brown("Could not detect your shell:"), err.Error())
}
