package features

import (
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register("golang", golangActivate, nil)
}

const (
	GolangSuffixMod    = "+mod"
	GolangSuffixGopath = "+gopath"
)

func golangActivate(ctx *context.Context, version string) (bool, error) {
	golang := helpers.NewGolang(ctx.Cfg, strings.Split(version, "+")[0])

	if !golang.Exists() {
		return true, nil
	}

	ctx.Env.PrependToPath(golang.BinPath())

	ctx.Env.Set("GOROOT", golang.Path())

	switch {
	case strings.HasSuffix(version, GolangSuffixMod):
		ctx.Env.Set("GO111MODULE", "on")
	case strings.HasSuffix(version, GolangSuffixGopath):
		ctx.Env.Set("GO111MODULE", "off")
	}

	return false, nil
}
