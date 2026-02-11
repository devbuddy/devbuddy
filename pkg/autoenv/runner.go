package autoenv

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv/features"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

// Sync activates / deactivates the features in the instance of env.Env.
// When a feature is already active but unknown, it will be ignored completely.
// When a param changes, the feature is deactivated with the current param then activated with the new param.
func Sync(ctx *context.Context, set FeatureSet) {
	state := &StateManager{ctx.Env, ctx.UI}
	runner := &runner{
		ctx:         ctx,
		state:       state,
		features:    features.All(),
		fileTracker: utils.NewFileTracker(state.GetFileChecksums()),
	}
	runner.sync(set)
}

type runner struct {
	ctx         *context.Context
	state       *StateManager
	features    features.Features
	fileTracker *utils.FileTracker
}

func (r *runner) sync(featureSet FeatureSet) {
	if r.state.GetProjectSlug() != "" {
		// A project was active until now

		if r.ctx.Project == nil {
			// We jumped out of the project

			r.state.RestoreEnv()
			r.state.ForgetEnv()
		} else if r.state.GetProjectSlug() != r.ctx.Project.Slug() {
			// We jumped to a different project

			// Since it's a different project, we just deactivate all features
			// For example, "python:3.7" is activating a virtualenv built for a specific project
			for _, featureInfo := range r.state.GetActiveFeatures() {
				r.deactivateFeature(featureInfo)
			}

			r.state.RestoreEnv()
			// No ForgetEnv(), we keep the SavedEnv until we jump out of a project
		}
	}

	activeFeatures := r.state.GetActiveFeatures()

	for _, name := range r.features.Names() {
		wantFeatureInfo := featureSet.Get(name)
		activeFeatureInfo := activeFeatures.Get(name)

		if wantFeatureInfo != nil {
			if activeFeatureInfo != nil {
				if wantFeatureInfo.Param != activeFeatureInfo.Param {
					r.deactivateFeature(activeFeatureInfo)
					r.activateFeature(wantFeatureInfo)
				} else {
					// Same param — check if watched files changed
					r.reactivateIfFilesChanged(wantFeatureInfo)
				}
			} else {
				r.activateFeature(wantFeatureInfo)
			}
		} else {
			if activeFeatureInfo != nil {
				r.deactivateFeature(activeFeatureInfo)
			}
		}
	}

	r.state.SetFileChecksums(r.fileTracker.Checksums())

	if r.ctx.Project != nil {
		// Record the project and the environment mutations made by this project
		r.state.SetProjectSlug(r.ctx.Project.Slug())
		r.state.SaveEnv()
	} else {
		// Record that we are NOT in a project
		r.state.SetProjectSlug("")
	}
}

func (r *runner) activateFeature(featureInfo *FeatureInfo) {
	r.ctx.UI.Debug("activating %s (%s)", featureInfo.Name, featureInfo.Param)

	environment := r.features.Get(featureInfo.Name)
	if environment == nil {
		r.ctx.UI.Warningf("unknown feature: %s (ignoring)", featureInfo.Name)
		return
	}

	devUpNeeded, err := environment.Activate(r.ctx, featureInfo.Param)
	if err != nil {
		r.ctx.UI.Debug("failed: %s", err)
		return
	}
	if devUpNeeded {
		r.ctx.UI.HookFeatureFailure(featureInfo.Name, featureInfo.Param)
		return
	}
	r.ctx.UI.HookFeatureActivated(featureInfo.Name, featureInfo.Param)
	r.state.SetFeature(featureInfo)

	// Prime the file tracker so the next check has a baseline
	r.primeFileTracker(featureInfo)
}

func (r *runner) reactivateIfFilesChanged(featureInfo *FeatureInfo) {
	def := r.features.Get(featureInfo.Name)
	if def == nil {
		return
	}
	watcher, ok := def.(features.FileWatcher)
	if !ok {
		return
	}
	for _, path := range watcher.WatchedFiles(featureInfo.Param) {
		if changed, _ := r.fileTracker.HasChanged(path); changed {
			r.ctx.UI.Debug("watched file changed: %s", path)
			r.activateFeature(featureInfo)
			return
		}
	}
}

func (r *runner) primeFileTracker(featureInfo *FeatureInfo) {
	def := r.features.Get(featureInfo.Name)
	if def == nil {
		return
	}
	watcher, ok := def.(features.FileWatcher)
	if !ok {
		return
	}
	for _, path := range watcher.WatchedFiles(featureInfo.Param) {
		_, _ = r.fileTracker.HasChanged(path) // stores the current checksum
	}
}

func (r *runner) deactivateFeature(featureInfo *FeatureInfo) {
	r.ctx.UI.Debug("deactivating %s (%s)", featureInfo.Name, featureInfo.Param)

	environment := r.features.Get(featureInfo.Name)
	if environment == nil {
		r.ctx.UI.Warningf("unknown feature: %s (ignoring)", featureInfo.Name)
		return
	}

	environment.Deactivate(r.ctx, featureInfo.Param)
	r.state.UnsetFeature(featureInfo.Name)
}
