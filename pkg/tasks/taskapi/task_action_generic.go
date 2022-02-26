package taskapi

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
)

type taskAction struct {
	desc       string
	conditions []Condition
	runFunc    func(*context.Context) error
	feature    *autoenv.FeatureInfo

	runCalled bool
}

func (a *taskAction) Description() string {
	return a.desc
}

func (a *taskAction) Needed(ctx *context.Context) (result *ActionResult) {
	if a.runCalled {
		return a.after(ctx)
	}
	return a.before(ctx)
}

func (a *taskAction) Run(ctx *context.Context) error {
	a.runCalled = true
	return a.runFunc(ctx)
}

func (a *taskAction) Feature() *autoenv.FeatureInfo {
	return a.feature
}

// internals

func (a *taskAction) before(ctx *context.Context) (result *ActionResult) {
	if len(a.conditions) == 0 {
		return Needed("action without conditions")
	}

	for _, condition := range a.conditions {
		result = condition.Before(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}

	return result
}

func (a *taskAction) after(ctx *context.Context) (result *ActionResult) {
	for _, condition := range a.conditions {
		result = condition.After(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	return NotNeeded()
}
