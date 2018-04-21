package features

import (
	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/env"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

type Runner struct {
	cfg  *config.Config
	proj *project.Project
	ui   *termui.HookUI
	env  *env.Env
}

func NewRunner(cfg *config.Config, proj *project.Project, ui *termui.HookUI, env *env.Env) *Runner {
	return &Runner{cfg: cfg, proj: proj, ui: ui, env: env}
}

func (r *Runner) Run(features map[string]string) {
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

func (r *Runner) activateFeature(name string, version string) {
	feature := allFeatures[name](version)

	err := feature.activate(r.cfg, r.proj, r.env)
	if err != nil {
		if err == DevUpNeeded {
			r.ui.HookFeatureFailure(name, version)
		} else {
			r.ui.Debug("failed: %s", err)
		}
	} else {
		r.ui.HookFeatureActivated(name, version)
		r.env.SetFeature(name, version)
	}
}

func (r *Runner) deactivateFeature(name string, version string) {
	feature := allFeatures[name](version)

	feature.deactivate(r.cfg, r.env)
	r.env.UnsetFeature(name)
	r.ui.Debug("%s deactivated", name)
}
