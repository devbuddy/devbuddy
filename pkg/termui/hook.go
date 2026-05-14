package termui

import (
	"fmt"
	"os"
	"strings"

	color "github.com/logrusorgru/aurora"
)

type Feature interface {
	FeatureName() string
	FeatureParam() string
}

func (u *UI) HookFeaturesActivated(features []Feature) {
	parts := make([]string, len(features))
	for i, f := range features {
		param := f.FeatureParam()
		if param == "" || strings.HasPrefix(param, "{") {
			parts[i] = fmt.Sprintf("%s", color.Blue(f.FeatureName()))
		} else {
			parts[i] = fmt.Sprintf("%s%s%s%s", color.Blue(f.FeatureName()), color.Gray(12, "["), color.Cyan(param), color.Gray(12, "]"))
		}
	}
	Fprintf(u.out, "🐼  %s %s\n", color.Cyan("activated:"), strings.Join(parts, color.Gray(12, ", ").String()))
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
