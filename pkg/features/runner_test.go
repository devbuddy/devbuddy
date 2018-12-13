package features

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/termui"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type recorder struct {
	entries []string
}

func (r *recorder) record(s ...string) {
	r.entries = append(r.entries, strings.Join(s, " "))
}

func (r *recorder) assert(t *testing.T, s ...string) {
	require.Equal(t, s, r.entries)
}

func (r *recorder) reset() {
	r.entries = []string{}
}

var featureCalls *recorder

func init() {
	featureCalls = &recorder{}

	rust := definitions.Register("rust")
	rust.Activate = func(param string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
		featureCalls.record("activate", "rust", param)
		return false, nil
	}
	rust.Deactivate = func(param string, cfg *config.Config, env *env.Env) {
		featureCalls.record("deactivate", "rust", param)
	}

	elixir := definitions.Register("elixir")
	elixir.Activate = func(param string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
		featureCalls.record("activate", "elixir", param)
		return false, nil
	}
	elixir.Deactivate = func(param string, cfg *config.Config, env *env.Env) {
		featureCalls.record("deactivate", "elixir", param)
	}
}

func TestRunner(t *testing.T) {
	_, ui := termui.NewTesting(false)

	runner := &Runner{
		cfg:  nil,
		proj: nil,
		ui:   ui,
		env:  env.New([]string{}),
	}

	wantedFeatures := map[string]string{}
	runner.Run(wantedFeatures)

	featureCalls.assert(t)

	wantedFeatures["rust"] = "1.0"
	runner.Run(wantedFeatures)
	featureCalls.assert(t, "activate rust 1.0")

	featureCalls.reset()

	wantedFeatures["rust"] = "2.0"
	runner.Run(wantedFeatures)
	featureCalls.assert(t, "deactivate rust 1.0", "activate rust 2.0")

	featureCalls.reset()

	delete(wantedFeatures, "rust")
	runner.Run(wantedFeatures)
	featureCalls.assert(t, "deactivate rust 2.0")
}
