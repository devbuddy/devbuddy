package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

type runner struct {
	cfg   *config.Config
	proj  *project.Project
	ui    *termui.UI
	env   *env.Env
	state *FeatureState
	reg   *featureRegister
}

func (r *runner) sync(featureSet FeatureSet) {
	activeFeatures := r.state.GetActiveFeatures()

	for _, name := range r.reg.names() {
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
}

func (r *runner) activateFeature(featureInfo FeatureInfo) {
	r.ui.Debug("activating %s (%s)", featureInfo.Name, featureInfo.Param)

	environment, err := r.reg.get(featureInfo.Name)
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

	environment, err := r.reg.get(featureInfo.Name)
	if err != nil {
		r.ui.Warningf("%s (ignoring)", err)
		return
	}

	environment.Deactivate(featureInfo.Param, r.cfg, r.env)
	r.state.UnsetFeature(featureInfo.Name)
}
