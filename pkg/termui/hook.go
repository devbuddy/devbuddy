package termui

import baseui "github.com/devbuddy/devbuddy/pkg/ui"

type Feature = baseui.Feature

func HookShellDetectionError(err error) {
	baseui.HookShellDetectionError(err)
}
