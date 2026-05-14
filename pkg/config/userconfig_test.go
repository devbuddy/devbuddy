package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadUserConfig_MissingFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg := LoadUserConfig()

	require.Equal(t, false, cfg.Shell.DeferInit)
}

func TestLoadUserConfig_ValidFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	configDir := filepath.Join(dir, "devbuddy")
	require.NoError(t, os.MkdirAll(configDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(configDir, "config.yml"),
		[]byte("shell:\n  defer_init: true\n"),
		0o644,
	))

	cfg := LoadUserConfig()

	require.Equal(t, true, cfg.Shell.DeferInit)
}

func TestLoadUserConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	configDir := filepath.Join(dir, "devbuddy")
	require.NoError(t, os.MkdirAll(configDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(configDir, "config.yml"),
		[]byte("not: [valid: yaml"),
		0o644,
	))

	cfg := LoadUserConfig()

	require.Equal(t, false, cfg.Shell.DeferInit)
}

func TestLoadUserConfig_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	configDir := filepath.Join(dir, "devbuddy")
	require.NoError(t, os.MkdirAll(configDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(configDir, "config.yml"),
		[]byte(""),
		0o644,
	))

	cfg := LoadUserConfig()

	require.Equal(t, false, cfg.Shell.DeferInit)
}
