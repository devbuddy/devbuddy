package taskapi

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv"
)

type taskActionBuilder struct {
	*taskAction
}

// On registers a new condition
func (a *taskActionBuilder) On(condition Condition) *taskActionBuilder {
	a.conditions = append(a.conditions, condition)
	return a
}

// SetFeature defines that the feature specified should be activated.
func (a *taskActionBuilder) SetFeature(name, param string) *taskActionBuilder {
	a.feature = autoenv.NewFeatureInfo(name, param)
	return a
}
