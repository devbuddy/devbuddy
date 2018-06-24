package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/pkg/utils"
	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"
)

func touchNow(t *testing.T, path string) {
	now := time.Now()
	require.NoError(t, utils.Touch(path, now, now))
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
	touchNow(t, path)

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
	touchNow(t, path)

	err := s.RecordFileChange("testfile")
	require.NoError(t, err)

	require.False(t, s.HasFileChanged("testfile"))

	touchNow(t, path)
	require.True(t, s.HasFileChanged("testfile"))
}
