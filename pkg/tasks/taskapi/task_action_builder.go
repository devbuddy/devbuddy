package taskapi

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
)

type genericTaskActionBuilder struct {
	*genericTaskAction
}

// On registers a new condition
func (a *genericTaskActionBuilder) On(condition *genericTaskActionCondition) *genericTaskActionBuilder {
	a.conditions = append(a.conditions, condition)
	return a
}

// OnFunc registers a condition defined as a single function
func (a *genericTaskActionBuilder) OnFunc(condFunc func(*context.Context) *ActionResult) *genericTaskActionBuilder {
	a.On(&genericTaskActionCondition{pre: condFunc, post: condFunc})
	return a
}

// OnFileChange specifies that the action will run when a file changes.
// The action will NOT run if the file does not exist.
// The action will NOT fail if the file is not created.
func (a *genericTaskActionBuilder) OnFileChange(path string) *genericTaskActionBuilder {
	a.monitoredFiles = append(a.monitoredFiles, path)
	return a
}

// SetFeature defines that the feature specified should be activated.
func (a *genericTaskActionBuilder) SetFeature(name, param string) *genericTaskActionBuilder {
	a.feature = autoenv.NewFeatureInfo(name, param)
	return a
}
