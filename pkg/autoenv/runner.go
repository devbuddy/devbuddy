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
		ctx:   ctx,
		state: &FeatureState{ctx.Env, ctx.UI},
		reg:   features.GlobalRegister(),
	}
	runner.sync(set)
}

type runner struct {
	ctx   *context.Context
	state *FeatureState
	reg   features.Register
}

func (r *runner) sync(featureSet FeatureSet) {
	if r.state.GetProjectSlug() != "" {
		// A project was active until now

		if r.ctx.Project == nil {
			// We jumped out of the project

			r.state.RestoreEnv()
			r.state.ForgetEnv() // we won't need to restore it in the future
		}

		if r.ctx.Project != nil && r.state.GetProjectSlug() != r.ctx.Project.Slug() {
			// We jumped to a different project

			r.state.RestoreEnv()
			// When we jump to another project, we will need to restore the initial environment eventually

			for _, featureInfo := range r.state.GetActiveFeatures() {
				r.deactivateFeature(featureInfo)
			}
		}
	}

	activeFeatures := r.state.GetActiveFeatures()

	for _, name := range r.reg.Names() {
		wantFeatureInfo := featureSet.Get(name)
		activeFeatureInfo := activeFeatures.Get(name)

		if wantFeatureInfo != nil {
			if activeFeatureInfo != nil {
				if wantFeatureInfo.Param != activeFeatureInfo.Param {
					r.deactivateFeature(activeFeatureInfo)
					r.activateFeature(wantFeatureInfo)
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

	// Record for which project the features were activated
	if r.ctx.Project != nil {
		r.state.SetProjectSlug(r.ctx.Project.Slug())
		err := r.state.SaveEnv()
		if err != nil {
			r.ctx.UI.Debug("state.RecordPrevious() failed with: %s", err)
		}
	} else {
		r.state.SetProjectSlug("")
	}
}

func (r *runner) activateFeature(featureInfo *FeatureInfo) {
	r.ctx.UI.Debug("activating %s (%s)", featureInfo.Name, featureInfo.Param)

	environment, err := r.reg.Get(featureInfo.Name)
	if err != nil {
		r.ctx.UI.Warningf("%s (ignoring)", err)
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
}

func (r *runner) deactivateFeature(featureInfo *FeatureInfo) {
	r.ctx.UI.Debug("deactivating %s (%s)", featureInfo.Name, featureInfo.Param)

	environment, err := r.reg.Get(featureInfo.Name)
	if err != nil {
		r.ctx.UI.Warningf("%s (ignoring)", err)
		return
	}

	environment.Deactivate(r.ctx, featureInfo.Param)
	r.state.UnsetFeature(featureInfo.Name)
}
