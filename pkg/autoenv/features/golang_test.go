package features

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/context"
	environ "github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/stretchr/testify/require"
)

func TestGolangActivate_AddsGoInstallBinFromGOBIN(t *testing.T) {
	ctx := goFeatureTestContext(t, []string{"PATH=/usr/bin", "GOBIN=/opt/go-tools/bin"})

	again, err := golang{}.Activate(ctx, "1.23.6")

	require.NoError(t, err)
	require.False(t, again)
	require.Equal(t, filepath.Join(goFeaturePath(t, ctx, "1.23.6"), "bin")+":/opt/go-tools/bin:/usr/bin", ctx.Env.Get("PATH"))
}

func TestGolangActivate_AddsGoInstallBinsFromGOPATH(t *testing.T) {
	ctx := goFeatureTestContext(t, []string{"PATH=/usr/bin", "GOPATH=/opt/go:/workspace/go"})

	again, err := golang{}.Activate(ctx, "1.23.6")

	require.NoError(t, err)
	require.False(t, again)
	require.Equal(t, filepath.Join(goFeaturePath(t, ctx, "1.23.6"), "bin")+":/workspace/go/bin:/opt/go/bin:/usr/bin", ctx.Env.Get("PATH"))
}

func TestGolangActivate_AddsDefaultGoInstallBin(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	ctx := goFeatureTestContext(t, []string{"PATH=/usr/bin"})

	again, err := golang{}.Activate(ctx, "1.23.6")

	require.NoError(t, err)
	require.False(t, again)
	require.Equal(t, filepath.Join(goFeaturePath(t, ctx, "1.23.6"), "bin")+":"+filepath.Join(home, "go", "bin")+":/usr/bin", ctx.Env.Get("PATH"))
}

func goFeatureTestContext(t *testing.T, envVars []string) *context.Context {
	t.Helper()

	t.Setenv("XDG_DATA_HOME", t.TempDir())
	cfg, err := config.Load()
	require.NoError(t, err)
	ctx := &context.Context{
		Cfg: cfg,
		Env: environ.New(envVars),
	}
	goPath := goFeaturePath(t, ctx, "1.23.6")
	require.NoError(t, os.MkdirAll(filepath.Join(goPath, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(goPath, "bin", "go"), []byte(""), 0o755))
	return ctx
}

func goFeaturePath(t *testing.T, ctx *context.Context, version string) string {
	t.Helper()
	return helpers.NewGolang(ctx, version).Path()
}
