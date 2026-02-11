package utils

import (
	"path/filepath"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestFileTracker_FirstCheck(t *testing.T) {
	_, tmpfile := test.CreateFile(t, ".env", []byte("VAR=one"))

	ft := NewFileTracker(nil)

	changed, err := ft.HasChanged(tmpfile)
	require.NoError(t, err)
	require.False(t, changed, "first check should return false (no baseline)")

	require.Contains(t, ft.Checksums(), tmpfile)
}

func TestFileTracker_Unchanged(t *testing.T) {
	_, tmpfile := test.CreateFile(t, ".env", []byte("VAR=one"))

	ft := NewFileTracker(nil)

	ft.HasChanged(tmpfile) // prime

	changed, err := ft.HasChanged(tmpfile)
	require.NoError(t, err)
	require.False(t, changed, "same content should not report changed")
}

func TestFileTracker_Changed(t *testing.T) {
	_, tmpfile := test.CreateFile(t, ".env", []byte("VAR=one"))

	ft := NewFileTracker(nil)

	ft.HasChanged(tmpfile) // prime

	test.WriteFile(tmpfile, []byte("VAR=two"))

	changed, err := ft.HasChanged(tmpfile)
	require.NoError(t, err)
	require.True(t, changed, "different content should report changed")
}

func TestFileTracker_MissingFile(t *testing.T) {
	tmpdir := t.TempDir()
	missing := filepath.Join(tmpdir, "nope.env")

	ft := NewFileTracker(nil)

	changed, err := ft.HasChanged(missing)
	require.NoError(t, err)
	require.False(t, changed, "missing file should return false")
}

func TestFileTracker_RestoredFromChecksums(t *testing.T) {
	_, tmpfile := test.CreateFile(t, ".env", []byte("VAR=one"))

	// Simulate a previous invocation that stored the checksum
	ft1 := NewFileTracker(nil)
	ft1.HasChanged(tmpfile)
	saved := ft1.Checksums()

	// New tracker initialized from persisted checksums
	ft2 := NewFileTracker(saved)

	// Same content → not changed
	changed, err := ft2.HasChanged(tmpfile)
	require.NoError(t, err)
	require.False(t, changed)

	// Modify the file → changed
	test.WriteFile(tmpfile, []byte("VAR=two"))

	changed, err = ft2.HasChanged(tmpfile)
	require.NoError(t, err)
	require.True(t, changed)
}
