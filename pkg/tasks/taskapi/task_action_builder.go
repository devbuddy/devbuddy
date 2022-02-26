package taskapi

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
)

type taskActionBuilder struct {
	*taskAction
}

// On registers a new condition
func (a *taskActionBuilder) On(condition *taskActionCondition) *taskActionBuilder {
	a.conditions = append(a.conditions, condition)
	return a
}

// OnFunc registers a condition defined as a single function
func (a *taskActionBuilder) OnFunc(condFunc func(*context.Context) *ActionResult) *taskActionBuilder {
	a.On(&taskActionCondition{pre: condFunc, post: condFunc})
	return a
}

// OnFileChange specifies that the action will run when a file changes.
// The action will NOT run if the file does not exist.
// The action will NOT fail if the file is not created.
func (a *taskActionBuilder) OnFileChange(path string) *taskActionBuilder {
	a.monitoredFiles = append(a.monitoredFiles, path)
	return a
}

// SetFeature defines that the feature specified should be activated.
func (a *taskActionBuilder) SetFeature(name, param string) *taskActionBuilder {
	a.feature = autoenv.NewFeatureInfo(name, param)
	return a
}
