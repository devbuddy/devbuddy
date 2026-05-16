package features

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/context"
	envpkg "github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/ui"

	"github.com/stretchr/testify/require"
)

func newRubyTestContext(t *testing.T, projectPath string) (*context.Context, func() string) {
	t.Helper()
	buf, ui := ui.NewBufferedTesting(false)
	ctx := &context.Context{
		Env:     envpkg.New([]string{}),
		UI:      ui,
		Project: project.NewFromPath(projectPath),
	}
	return ctx, buf.String
}

func TestRubyWarnsOnRubyVersionMismatch(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ruby-version"), []byte("3.2.0\n"), 0o600))

	ctx, output := newRubyTestContext(t, dir)
	_, _ = ruby{}.Activate(ctx, "3.3.0")

	require.Contains(t, output(), "dev.yml requests 3.3.0 but .ruby-version says 3.2.0")
}

func TestRubySilentWhenRubyVersionMatches(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ruby-version"), []byte("3.3.0\n"), 0o600))

	ctx, output := newRubyTestContext(t, dir)
	_, _ = ruby{}.Activate(ctx, "3.3.0")

	require.NotContains(t, output(), ".ruby-version")
}

func TestRubySilentWhenRubyVersionAbsent(t *testing.T) {
	dir := t.TempDir()
	ctx, output := newRubyTestContext(t, dir)
	_, _ = ruby{}.Activate(ctx, "3.3.0")

	require.NotContains(t, output(), ".ruby-version")
}

func TestRubyWatchesRubyVersionFile(t *testing.T) {
	require.Equal(t, []string{".ruby-version"}, ruby{}.WatchedFiles("any"))
}
