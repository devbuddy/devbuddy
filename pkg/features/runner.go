package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
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

	for _, name := range definitions.Names() {
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

func (r *Runner) activateFeature(name string, param string) {
	err := Activate(name, param, r.cfg, r.proj, r.env)
	if err != nil {
		if err == DevUpNeeded {
			r.ui.HookFeatureFailure(name, param)
		} else {
			r.ui.Debug("failed: %s", err)
		}
	} else {
		r.ui.HookFeatureActivated(name, param)
		r.env.SetFeature(name, param)
	}
}

func (r *Runner) deactivateFeature(name string, param string) {
	Deactivate(name, param, r.cfg, r.env)
	r.env.UnsetFeature(name)
	r.ui.Debug("%s deactivated", name)
}
