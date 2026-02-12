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

	checksums, err := state.GetFileChecksums()
	if err != nil {
		ctx.UI.Warningf("autoenv: failed to read state: %s", err)
		return
	}

	runner := &runner{
		ctx:         ctx,
		state:       state,
		features:    features.All(),
		fileTracker: utils.NewFileTracker(checksums),
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
	projectSlug, err := r.state.GetProjectSlug()
	if err != nil {
		r.ctx.UI.Warningf("autoenv: failed to read state: %s", err)
		return
	}

	if projectSlug != "" {
		// A project was active until now

		if r.ctx.Project == nil {
			// We jumped out of the project

			if err := r.state.RestoreEnv(); err != nil {
				r.ctx.UI.Warningf("autoenv: failed to restore env: %s", err)
				return
			}
			if err := r.state.ForgetEnv(); err != nil {
				r.ctx.UI.Warningf("autoenv: failed to forget env: %s", err)
				return
			}
		} else if projectSlug != r.ctx.Project.Slug() {
			// We jumped to a different project

			// Since it's a different project, we just deactivate all features
			// For example, "python:3.7" is activating a virtualenv built for a specific project
			activeFeatures, err := r.state.GetActiveFeatures()
			if err != nil {
				r.ctx.UI.Warningf("autoenv: failed to read state: %s", err)
				return
			}
			for _, featureInfo := range activeFeatures {
				r.deactivateFeature(featureInfo)
			}

			if err := r.state.RestoreEnv(); err != nil {
				r.ctx.UI.Warningf("autoenv: failed to restore env: %s", err)
				return
			}
			// No ForgetEnv(), we keep the SavedEnv until we jump out of a project
		}
	}

	activeFeatures, err := r.state.GetActiveFeatures()
	if err != nil {
		r.ctx.UI.Warningf("autoenv: failed to read state: %s", err)
		return
	}

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

	if err := r.state.SetFileChecksums(r.fileTracker.Checksums()); err != nil {
		r.ctx.UI.Warningf("autoenv: failed to write state: %s", err)
		return
	}

	if r.ctx.Project != nil {
		// Record the project and the environment mutations made by this project
		if err := r.state.SetProjectSlug(r.ctx.Project.Slug()); err != nil {
			r.ctx.UI.Warningf("autoenv: failed to write state: %s", err)
			return
		}
		if err := r.state.SaveEnv(); err != nil {
			r.ctx.UI.Warningf("autoenv: failed to save env: %s", err)
			return
		}
	} else {
		// Record that we are NOT in a project
		if err := r.state.SetProjectSlug(""); err != nil {
			r.ctx.UI.Warningf("autoenv: failed to write state: %s", err)
			return
		}
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
	if err := r.state.SetFeature(featureInfo); err != nil {
		r.ctx.UI.Warningf("autoenv: failed to write state: %s", err)
		return
	}

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
	if err := r.state.UnsetFeature(featureInfo.Name); err != nil {
		r.ctx.UI.Warningf("autoenv: failed to write state: %s", err)
	}
}
