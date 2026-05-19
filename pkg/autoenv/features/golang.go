package features

import (
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(golang{})
}

type golang struct{}

func (golang) Name() string {
	return "golang"
}

func (golang) Activate(ctx *context.Context, param string) (bool, error) {
	golang := helpers.NewGolang(ctx, param)

	if !golang.Exists() {
		return true, nil
	}

	for _, binPath := range goInstallBinPaths(ctx) {
		ctx.Env.PrependToPath(binPath)
	}
	ctx.Env.PrependToPath(golang.BinPath())

	ctx.Env.Set("GOROOT", golang.Path())

	return false, nil
}

func (golang) Deactivate(ctx *context.Context, param string) {}

func goInstallBinPaths(ctx *context.Context) []string {
	if gobin := ctx.Env.Get("GOBIN"); gobin != "" {
		return []string{gobin}
	}

	if gopath := ctx.Env.Get("GOPATH"); gopath != "" {
		paths := []string{}
		for _, dir := range filepath.SplitList(gopath) {
			if dir == "" {
				continue
			}
			paths = append(paths, filepath.Join(dir, "bin"))
		}
		return paths
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{filepath.Join(home, "go", "bin")}
}
