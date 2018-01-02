package features

import (
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

type Runner struct {
	proj *project.Project
	ui   *termui.HookUI
	env  *Env
}

func NewRunner(proj *project.Project, ui *termui.HookUI, env *Env) *Runner {
	return &Runner{proj: proj, ui: ui, env: env}
}

func (r *Runner) Run() {
	wantedFeatures := r.getWantedFeatures()
	r.handleFeatures(wantedFeatures)
	r.env.SetActiveFeatures(wantedFeatures)
}

func (r *Runner) handleFeatures(features map[string]string) {
	activeFeatures := r.env.GetActiveFeatures()

	for name := range allFeatures {
		wantVersion, want := features[name]
		activeVersion, active := activeFeatures[name]

		if want {
			if !active || wantVersion != activeVersion {
				r.activateFeature(name, wantVersion)
			}
		} else {
			if active {
				r.deactivateFeature(name, activeVersion)
			}
		}
	}

}

func (r *Runner) getWantedFeatures() map[string]string {
	var err error
	wantedFeatures := map[string]string{}

	if r.proj != nil {
		wantedFeatures, err = r.proj.GetFeatures()
		if err != nil {
			r.ui.Debug("failed to get project tasks: %s", err)
		}
	}
	return wantedFeatures
}

func (r *Runner) activateFeature(name string, version string) {
	feature := allFeatures[name](version)

	err := feature.Enable(r.proj, r.env, r.ui)
	if err != nil {
		if err == DevUpNeeded {
			r.ui.HookFeatureFailure(name, version)
		} else {
			r.ui.Debug("failed: %s", err)
		}
	} else {
		r.ui.HookFeatureActivated(name, version)
	}
}

func (r *Runner) deactivateFeature(name string, version string) {
	feature := allFeatures[name](version)

	feature.Disable(r.proj, r.env, r.ui)
	r.ui.Debug("%s deactivated", name)
}
