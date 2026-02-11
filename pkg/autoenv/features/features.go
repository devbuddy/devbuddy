package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
)

type Feature interface {
	Name() string
	Activate(ctx *context.Context, param string) (devUpNeeded bool, err error)
	Deactivate(ctx *context.Context, param string)
}

// FileWatcher is an optional interface for features that depend on external files.
// When a watched file changes on disk, the feature will be re-activated.
type FileWatcher interface {
	WatchedFiles(param string) []string
}

type Features interface {
	Get(string) Feature
	Names() []string
}

func All() Features {
	return &register
}

var register = Register{}
