package autoenv

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv/features"
	"github.com/devbuddy/devbuddy/pkg/context"
)

// Sync activates / deactivates the features in the instance of env.Env.
// When a feature is already active but unknown, it will be ignored completely.
// When a param changes, the feature is deactivated with the current param then activated with the new param.
func Sync(ctx *context.Context, set FeatureSet) {
	runner := &runner{
		ctx:      ctx,
		state:    &StateManager{ctx.Env, ctx.UI},
		features: features.All(),
	}
	runner.sync(set)
}

type runner struct {
	ctx      *context.Context
	state    *StateManager
	features features.Features
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
			for _, fi := range r.state.GetActiveFeatures() {
				r.deactivate(fi)
			}

			r.state.RestoreEnv()
			// No ForgetEnv(), we keep the SavedEnv until we jump out of a project
		}
	}

	activeFeatures := r.state.GetActiveFeatures()

	for _, name := range r.features.Names() {
		want := featureSet.Get(name)
		active := activeFeatures.Get(name)

		if want != nil {
			if active != nil {
				if want.Param == active.Param {
					r.refresh(want)
				} else {
					r.deactivate(active)
					r.activate(want)
				}
			} else {
				r.activate(want)
			}
		} else {
			if active != nil {
				r.deactivate(active)
			}
		}
	}

	if r.ctx.Project != nil {
		// Record the project and the environment mutations made by this project
		r.state.SetProjectSlug(r.ctx.Project.Slug())
		r.state.SaveEnv()
	} else {
		// Record that we are NOT in a project
		r.state.SetProjectSlug("")
	}
}

func (r *runner) activate(fi *FeatureInfo) {
	r.ctx.UI.Debug("activating %s (%s)", fi.Name, fi.Param)

	feature := r.features.Get(fi.Name)
	if feature == nil {
		r.ctx.UI.Warningf("unknown feature: %s (ignoring)", fi.Name)
		return
	}

	devUpNeeded, err := feature.Activate(r.ctx, fi.Param)
	if err != nil {
		r.ctx.UI.Debug("failed: %s", err)
		return
	}
	if devUpNeeded {
		r.ctx.UI.HookFeatureFailure(fi.Name, fi.Param)
		return
	}
	r.ctx.UI.HookFeatureActivated(fi.Name, fi.Param)
	r.state.SetFeature(fi)
}

func (r *runner) refresh(fi *FeatureInfo) {
	r.ctx.UI.Debug("refresh %s (%s)", fi.Name, fi.Param)

	feature := r.features.Get(fi.Name)
	if feature == nil {
		r.ctx.UI.Warningf("unknown feature: %s (ignoring)", fi.Name)
		return
	}

	feature.Refresh(r.ctx, fi.Param)
}

func (r *runner) deactivate(fi *FeatureInfo) {
	r.ctx.UI.Debug("deactivating %s (%s)", fi.Name, fi.Param)

	feature := r.features.Get(fi.Name)
	if feature == nil {
		r.ctx.UI.Warningf("unknown feature: %s (ignoring)", fi.Name)
		return
	}

	feature.Deactivate(r.ctx, fi.Param)
	r.state.UnsetFeature(fi.Name)
}
