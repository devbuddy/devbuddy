package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

type Runner struct {
	cfg  *config.Config
	proj *project.Project
	ui   *termui.UI
	env  *env.Env
}

func NewRunner(cfg *config.Config, proj *project.Project, ui *termui.UI, env *env.Env) *Runner {
	return &Runner{cfg: cfg, proj: proj, ui: ui, env: env}
}

func (r *Runner) Run(features map[string]string) {
	activeFeatures := r.env.GetActiveFeatures()

	for _, name := range definitions.Names() {
		wantVersion, want := features[name]
		activeVersion, active := activeFeatures[name]

		if want {
			if active {
				if wantVersion != activeVersion {
					r.deactivateFeature(name, activeVersion)
					r.activateFeature(name, wantVersion)
				}
			} else {
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
	r.ui.Debug("activating %s (%s)", name, param)

	devUpNeeded, err := Activate(name, param, r.cfg, r.proj, r.env)
	if err != nil {
		r.ui.Debug("failed: %s", err)
		return
	}
	if devUpNeeded {
		r.ui.HookFeatureFailure(name, param)
		return
	}
	r.ui.HookFeatureActivated(name, param)
	r.env.SetFeature(name, param)
}

func (r *Runner) deactivateFeature(name string, param string) {
	r.ui.Debug("deactivating %s (%s)", name, param)

	Deactivate(name, param, r.cfg, r.env)
	r.env.UnsetFeature(name)
}
