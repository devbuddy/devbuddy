package store

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"
)

func TestWithoutFile(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	result, err := s.HasFileChanged("testfile")
	require.NoError(t, err)
	require.True(t, result)
}

func TestFirstTime(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	filet.File(t, filepath.Join(tmpdir, "testfile"), "some-value")

	result, err := s.HasFileChanged("testfile")
	require.NoError(t, err)
	require.True(t, result)
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

	filet.File(t, filepath.Join(tmpdir, "testfile"), "some-value")

	err := s.RecordFileChange("testfile")
	require.NoError(t, err)

	result, err := s.HasFileChanged("testfile")
	require.NoError(t, err)
	require.False(t, result)

	filet.File(t, filepath.Join(tmpdir, "testfile"), "some-OTHER-value")

	result, err = s.HasFileChanged("testfile")
	require.NoError(t, err)
	require.True(t, result)
}
