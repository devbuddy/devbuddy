package store

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"
)

func TestProjectPathMissing(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir + "/nopenope")

	err := s.Set("dummy", []byte(""))
	require.Error(t, err)
}

func TestInitialization(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	err := s.SetString("dummykey", "dummy")
	require.NoError(t, err)

	filet.DirContains(t, tmpdir, ".devbuddy")
	filet.DirContains(t, tmpdir, ".devbuddy/dummykey")

	filet.DirContains(t, tmpdir, ".devbuddy/.gitignore")
	filet.FileSays(t, tmpdir+"/.devbuddy/.gitignore", []byte("*"))
}

func TestSetGet(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	testValues := [][]byte{
		[]byte("DUMMY"),
		[]byte(""),
		[]byte("   "),
	}

	for _, testVal := range testValues {
		err := s.Set("key", testVal)
		require.NoError(t, err)

		val, err := s.Get("key")
		require.NoError(t, err)
		require.Equal(t, testVal, val)
	}
}

func TestSetGetString(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	testValues := []string{
		"DUMMY",
		"",
		"   ",
	}

	for _, testVal := range testValues {
		err := s.SetString("key", testVal)
		require.NoError(t, err)

		val, err := s.GetString("key")
		require.NoError(t, err)
		require.Equal(t, testVal, val)
	}
}

func TestKeyEmpty(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	_, err := s.Get("")
	require.Error(t, err)

	err = s.Set("", []byte(""))
	require.Error(t, err)
}

func TestGetNotFound(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	val, err := s.Get("nope")
	require.NoError(t, err)
	require.Equal(t, []byte(nil), val)
}

func TestGetStringNotFound(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	val, err := s.GetString("nope")
	require.NoError(t, err)
	require.Equal(t, "", val)
}

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
