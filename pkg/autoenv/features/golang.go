package features

import (
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(Golang{})
}

const (
	GolangSuffixMod    = "+mod"
	GolangSuffixGopath = "+gopath"
)

type Golang struct{}

func (Golang) Name() string {
	return "golang"
}

func (Golang) Activate(ctx *context.Context, param string) (bool, error) {
	golang := helpers.NewGolang(ctx.Cfg, strings.Split(param, "+")[0])

	if !golang.Exists() {
		return true, nil
	}

	ctx.Env.PrependToPath(golang.BinPath())

	ctx.Env.Set("GOROOT", golang.Path())

	switch {
	case strings.HasSuffix(param, GolangSuffixMod):
		ctx.Env.Set("GO111MODULE", "on")
	case strings.HasSuffix(param, GolangSuffixGopath):
		ctx.Env.Set("GO111MODULE", "off")
	}

	return false, nil
}

func (Golang) Deactivate(ctx *context.Context, param string) {}
