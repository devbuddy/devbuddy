package features

import (
	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/tasks"
	"github.com/pior/dad/pkg/termui"
)

type Runner struct {
	cfg  *config.Config
	proj *project.Project
	ui   *termui.HookUI
	env  *Env
}

func NewRunner(cfg *config.Config, proj *project.Project, ui *termui.HookUI, env *Env) *Runner {
	return &Runner{cfg: cfg, proj: proj, ui: ui, env: env}
}

func (r *Runner) Run() {
	wantedFeatures, err := r.getWantedFeatures()
	if err != nil {
		r.ui.Debug("failed to get project tasks: %s", err)
	}
	r.ui.Debug("DEV_AUTO_ENV_FEATURES=\"%s\"", r.env.Get("DEV_AUTO_ENV_FEATURES"))
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

func (r *Runner) getWantedFeatures() (map[string]string, error) {
	wantedFeatures := map[string]string{}

	if r.proj != nil {
		taskList, err := tasks.GetTasksFromProject(r.proj)
		if err != nil {
			return nil, err
		}

		for _, task := range taskList {
			if t, ok := task.(tasks.TaskWithFeature); ok {
				feature, param := t.Feature(r.proj)
				wantedFeatures[feature] = param
			}
		}
	}

	return wantedFeatures, nil
}

func (r *Runner) activateFeature(name string, version string) {
	feature := allFeatures[name](version)

	err := feature.Enable(r.cfg, r.proj, r.env, r.ui)
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

	feature.Disable(r.cfg, r.env, r.ui)
	r.ui.Debug("%s deactivated", name)
}
