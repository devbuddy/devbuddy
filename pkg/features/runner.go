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
	ui   *termui.UI
	env  *env.Env
	reg  *featureRegister
}

func Sync(cfg *config.Config, proj *project.Project, ui *termui.UI, env *env.Env, features map[string]string) {
	runner := &Runner{cfg: cfg, proj: proj, ui: ui, env: env, reg: globalRegister}
	runner.sync(features)
}

func (r *Runner) sync(features map[string]string) {
	activeFeatures := r.env.GetActiveFeatures()

	for _, name := range r.reg.names() {
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

	environment, err := r.reg.get(name)
	if err != nil {
		r.ui.Warningf("%s (ignoring)", err)
		return
	}

	devUpNeeded, err := environment.Activate(param, r.cfg, r.proj, r.env)
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

	environment, err := r.reg.get(name)
	if err != nil {
		r.ui.Warningf("%s (ignoring)", err)
		return
	}

	environment.Deactivate(param, r.cfg, r.env)
	r.env.UnsetFeature(name)
}
