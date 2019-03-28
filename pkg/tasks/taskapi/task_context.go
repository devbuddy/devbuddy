package taskapi

import (
	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

// Context provides the context in which a Task run
type Context struct {
	Project  *project.Project
	UI       *termui.UI
	Cfg      *config.Config
	Env      *env.Env
	Features autoenv.FeatureSet
}
