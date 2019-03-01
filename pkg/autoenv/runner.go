package autoenv

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv/register"
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

// Sync activates / deactivates the features in the instance of env.Env.
// When a feature is already active but unknown, it will be ignored completely.
// When a param changes, the feature is deactivated with the current param then activated with the new param.
func Sync(cfg *config.Config, proj *project.Project, ui *termui.UI, env *env.Env, set FeatureSet) {
	runner := &runner{cfg: cfg, proj: proj, ui: ui, env: env, state: &FeatureState{env}, reg: register.Global()}
	runner.sync(set)
}

type runner struct {
	cfg   *config.Config
	proj  *project.Project
	ui    *termui.UI
	env   *env.Env
	state *FeatureState
	reg   register.ImmutableRegister
}

func (r *runner) sync(featureSet FeatureSet) {
	// If we jumped to a different project, all feature should be deactivated
	if r.proj != nil && r.state.GetProjectSlug() != r.proj.Slug() {
		for _, featureInfo := range r.state.GetActiveFeatures() {
			r.deactivateFeature(featureInfo)
		}
	}

	activeFeatures := r.state.GetActiveFeatures()

	for _, name := range r.reg.Names() {
		wantFeatureInfo, want := featureSet[name]
		activeFeatureInfo, active := activeFeatures[name]

		if want {
			if active {
				if wantFeatureInfo.Param != activeFeatureInfo.Param {
					r.deactivateFeature(activeFeatureInfo)
					r.activateFeature(wantFeatureInfo)
				}
			} else {
				r.activateFeature(wantFeatureInfo)
			}
		} else {
			if active {
				r.deactivateFeature(activeFeatureInfo)
			}
		}
	}

	// Record for which project the features were activated
	if r.proj != nil {
		r.state.SetProjectSlug(r.proj.Slug())
	}
}

func (r *runner) activateFeature(featureInfo FeatureInfo) {
	r.ui.Debug("activating %s (%s)", featureInfo.Name, featureInfo.Param)

	environment, err := r.reg.Get(featureInfo.Name)
	if err != nil {
		r.ui.Warningf("%s (ignoring)", err)
		return
	}

	devUpNeeded, err := environment.Activate(featureInfo.Param, r.cfg, r.proj, r.env)
	if err != nil {
		r.ui.Debug("failed: %s", err)
		return
	}
	if devUpNeeded {
		r.ui.HookFeatureFailure(featureInfo.Name, featureInfo.Param)
		return
	}
	r.ui.HookFeatureActivated(featureInfo.Name, featureInfo.Param)
	r.state.SetFeature(featureInfo)
}

func (r *runner) deactivateFeature(featureInfo FeatureInfo) {
	r.ui.Debug("deactivating %s (%s)", featureInfo.Name, featureInfo.Param)

	environment, err := r.reg.Get(featureInfo.Name)
	if err != nil {
		r.ui.Warningf("%s (ignoring)", err)
		return
	}

	environment.Deactivate(featureInfo.Param, r.cfg, r.env)
	r.state.UnsetFeature(featureInfo.Name)
}
