package termui

import (
	"os"
	"strings"

	baseui "github.com/devbuddy/devbuddy/pkg/ui"
)

type Feature interface {
	FeatureName() string
	FeatureParam() string
}

func (u *UI) HookFeaturesActivated(features []Feature) {
	plainParts := make([]string, len(features))
	fields := make([]baseui.Field, len(features))
	for i, f := range features {
		fields[i] = baseui.F(f.FeatureName(), f.FeatureParam())
		param := f.FeatureParam()
		if param == "" || strings.HasPrefix(param, "{") {
			plainParts[i] = f.FeatureName()
		} else {
			plainParts[i] = f.FeatureName() + "[" + param + "]"
		}
	}
	u.emit(baseui.Event{Kind: baseui.KindHookActivated, Text: strings.Join(plainParts, ", "), Fields: fields})
}

func (u *UI) HookFeatureFailure(name string, param string) {
	u.emit(baseui.Event{Kind: baseui.KindHookFeatureFailed, Text: name, Fields: []baseui.Field{baseui.F("param", param)}})
}

func (u *UI) HookDevYmlChanged() {
	u.emit(baseui.Event{Kind: baseui.KindHookDevYMLChanged})
}

func HookShellDetectionError(err error) {
	Fprintf(os.Stderr, "%s", baseui.TerminalRenderer{}.Render(baseui.Event{Kind: baseui.KindShellDetectError, Text: err.Error()}))
}
