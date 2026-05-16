package api

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/ui"
)

// Context provides the context in which a Task run
type Context struct {
	Project  *project.Project
	UI       *ui.UI
	Cfg      *config.Config
	Env      *env.Env
	Features autoenv.FeatureSet
}
