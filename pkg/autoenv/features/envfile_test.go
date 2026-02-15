package features

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeEnvFile(t *testing.T, dir string, content string) string {
	t.Helper()
	path := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestEnvfileActivate(t *testing.T) {
	ctx := newTestContext()
	dir := t.TempDir()
	path := writeEnvFile(t, dir, "A=1\nB=2\n")

	devUpNeeded, err := envfile{}.Activate(ctx, path)
	require.NoError(t, err)
	require.False(t, devUpNeeded)

	require.Equal(t, "1", ctx.Env.Get("A"))
	require.Equal(t, "2", ctx.Env.Get("B"))
	require.Equal(t, "A,B", ctx.Env.Get(envfileTrackedVarsKey))
}

func TestEnvfileActivateUnsetsRemovedVars(t *testing.T) {
	ctx := newTestContext()
	dir := t.TempDir()
	path := writeEnvFile(t, dir, "A=1\nB=2\n")

	// First activation: sets A and B
	_, err := envfile{}.Activate(ctx, path)
	require.NoError(t, err)
	require.Equal(t, "1", ctx.Env.Get("A"))
	require.Equal(t, "2", ctx.Env.Get("B"))

	// Modify .env: remove B
	writeEnvFile(t, dir, "A=1\n")

	// Second activation: should unset B
	_, err = envfile{}.Activate(ctx, path)
	require.NoError(t, err)

	require.Equal(t, "1", ctx.Env.Get("A"))
	require.Equal(t, "", ctx.Env.Get("B"))
	require.False(t, ctx.Env.Has("B"))
	require.Equal(t, "A", ctx.Env.Get(envfileTrackedVarsKey))
}

func TestEnvfileActivateUpdatesChangedValues(t *testing.T) {
	ctx := newTestContext()
	dir := t.TempDir()
	path := writeEnvFile(t, dir, "A=1\n")

	_, err := envfile{}.Activate(ctx, path)
	require.NoError(t, err)
	require.Equal(t, "1", ctx.Env.Get("A"))

	writeEnvFile(t, dir, "A=updated\n")

	_, err = envfile{}.Activate(ctx, path)
	require.NoError(t, err)
	require.Equal(t, "updated", ctx.Env.Get("A"))
}

func TestEnvfileActivateNoFile(t *testing.T) {
	ctx := newTestContext()

	devUpNeeded, err := envfile{}.Activate(ctx, "/nonexistent/.env")
	require.Error(t, err)
	require.True(t, devUpNeeded)
}

func TestEnvfileDeactivate(t *testing.T) {
	ctx := newTestContext()
	dir := t.TempDir()
	path := writeEnvFile(t, dir, "A=1\nB=2\n")

	_, err := envfile{}.Activate(ctx, path)
	require.NoError(t, err)
	require.Equal(t, "1", ctx.Env.Get("A"))
	require.Equal(t, "2", ctx.Env.Get("B"))

	envfile{}.Deactivate(ctx, path)

	require.False(t, ctx.Env.Has("A"))
	require.False(t, ctx.Env.Has("B"))
	require.False(t, ctx.Env.Has(envfileTrackedVarsKey))
}

func TestEnvfileDeactivateNoopWhenNotActivated(t *testing.T) {
	ctx := newTestContext()

	// Should not panic or error
	envfile{}.Deactivate(ctx, ".env")
}
