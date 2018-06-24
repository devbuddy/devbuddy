package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"
)

func touch(t *testing.T, path string) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	now := time.Now()
	err = os.Chtimes(path, now, now)
	require.NoError(t, err)
}

func TestWithoutFile(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	require.True(t, s.HasFileChanged("testfile"))
}

func TestFirstTime(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	path := filepath.Join(tmpdir, "testfile")
	touch(t, path)

	require.True(t, s.HasFileChanged("testfile"))
}

func TestRecordWithoutFile(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	err := s.RecordFileChange("testfile")
	require.Error(t, err)
}

func TestRecord(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	path := filepath.Join(tmpdir, "testfile")
	touch(t, path)

	err := s.RecordFileChange("testfile")
	require.NoError(t, err)

	require.False(t, s.HasFileChanged("testfile"))

	// time.Sleep(100 * time.Millisecond)
	touch(t, path)
	require.True(t, s.HasFileChanged("testfile"))
}
