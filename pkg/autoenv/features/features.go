package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
)

type Feature interface {
	Name() string
	Activate(ctx *context.Context, param string) (devUpNeeded bool, err error)
	Deactivate(ctx *context.Context, param string)
}

type Features interface {
	Get(string) Feature
	Names() []string
}

func All() Features {
	return &register
}

var register = Register{}
