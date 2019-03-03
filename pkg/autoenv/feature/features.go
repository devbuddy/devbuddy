package feature

import (
	"github.com/devbuddy/devbuddy/pkg/context"
)

type ActivateFunc func(*context.Context, string) (bool, error)
type DeactivateFunc func(*context.Context, string)

// Feature is the implementation of an environment feature.
type Feature struct {
	Name       string
	Activate   ActivateFunc
	Deactivate DeactivateFunc
}
