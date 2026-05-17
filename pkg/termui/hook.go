package termui

import (
	"fmt"
	"os"
	"strings"

	color "github.com/logrusorgru/aurora"

	baseui "github.com/devbuddy/devbuddy/pkg/ui"
)

type Feature interface {
	FeatureName() string
	FeatureParam() string
}

func (u *UI) HookFeaturesActivated(features []Feature) {
	parts := make([]string, len(features))
	plainParts := make([]string, len(features))
	fields := make([]baseui.Field, len(features))
	for i, f := range features {
		fields[i] = baseui.F(f.FeatureName(), f.FeatureParam())
		param := f.FeatureParam()
		if param == "" || strings.HasPrefix(param, "{") {
			plainParts[i] = f.FeatureName()
			parts[i] = fmt.Sprintf("%s", color.Blue(f.FeatureName()))
		} else {
			plainParts[i] = fmt.Sprintf("%s[%s]", f.FeatureName(), param)
			parts[i] = fmt.Sprintf("%s%s%s%s", color.Blue(f.FeatureName()), color.Gray(12, "["), color.Cyan(param), color.Gray(12, "]"))
		}
	}
	u.record(baseui.Event{Kind: baseui.KindHookActivated, Text: strings.Join(plainParts, ", "), Fields: fields})
	Fprintf(u.out, "🐼  %s %s\n", color.Cyan("activated:"), strings.Join(parts, color.Gray(12, ", ").String()))
}

func (u *UI) HookFeatureFailure(name string, param string) {
	u.record(baseui.Event{Kind: baseui.KindHookFeatureFailed, Text: name, Fields: []baseui.Field{baseui.F("param", param)}})
	msg := fmt.Sprintf("failed to activate %s. Try running 'bud up' first!", name)

	paramStr := ""
	if param != "" {
		paramStr = fmt.Sprintf(" (%s)", param)
	}

	Fprintf(u.out, "🐼  %s%s\n", color.Red(msg), color.Yellow(paramStr))
}

func (u *UI) HookDevYmlChanged() {
	u.record(baseui.Event{Kind: baseui.KindHookDevYMLChanged})
	Fprintf(u.out, "🐼  %s\n", color.Yellow("dev.yml changed, run `bud up` to apply"))
}

func HookShellDetectionError(err error) {
	Fprintf(os.Stderr, "%s %s\n", color.Yellow("Could not detect your shell:"), err.Error())
}
