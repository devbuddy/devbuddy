package termui

import (
	"fmt"
	"os"
	"strings"

	color "github.com/logrusorgru/aurora"
)

func (u *UI) HookFeaturesActivated(features []string) {
	Fprintf(u.out, "🐼  %s %s\n", color.Cyan("activated:"), color.Blue(strings.Join(features, ", ")))
}

func (u *UI) HookFeatureFailure(name string, param string) {
	msg := fmt.Sprintf("failed to activate %s. Try running 'bud up' first!", name)

	paramStr := ""
	if param != "" {
		paramStr = fmt.Sprintf(" (%s)", param)
	}

	Fprintf(u.out, "🐼  %s%s\n", color.Red(msg), color.Yellow(paramStr))
}

func (u *UI) HookDevYmlChanged() {
	Fprintf(u.out, "🐼  %s\n", color.Yellow("dev.yml changed, run `bud up` to apply"))
}

func HookShellDetectionError(err error) {
	Fprintf(os.Stderr, "%s %s\n", color.Yellow("Could not detect your shell:"), err.Error())
}
