package ui

import (
	"os"
	"strings"
)

type Feature interface {
	FeatureName() string
	FeatureParam() string
}

func (u *UI) HookFeaturesActivated(features []Feature) {
	plainParts := make([]string, len(features))
	fields := make([]Field, len(features))
	for i, f := range features {
		fields[i] = F(f.FeatureName(), f.FeatureParam())
		param := f.FeatureParam()
		if param == "" || strings.HasPrefix(param, "{") {
			plainParts[i] = f.FeatureName()
		} else {
			plainParts[i] = f.FeatureName() + "[" + param + "]"
		}
	}
	u.emit(Event{Kind: KindHookActivated, Text: strings.Join(plainParts, ", "), Fields: fields})
}

func (u *UI) HookFeatureFailure(name string, param string) {
	u.emit(Event{Kind: KindHookFeatureFailed, Text: name, Fields: []Field{F("param", param)}})
}

func (u *UI) HookDevYmlChanged() {
	u.emit(Event{Kind: KindHookDevYMLChanged})
}

func HookShellDetectionError(err error) {
	Fprintf(os.Stderr, "%s", TerminalRenderer{}.Render(Event{Kind: KindShellDetectError, Text: err.Error()}))
}
