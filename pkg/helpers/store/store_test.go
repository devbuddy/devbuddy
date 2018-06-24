package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"
)

func touch(t *testing.T, path string) {
	err := ioutil.WriteFile(path, []byte(""), 0644)
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

	time.Sleep(100 * time.Millisecond)
	touch(t, path)
	require.True(t, s.HasFileChanged("testfile"))
}

func TestMtimeProperties(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	test1 := filepath.Join(tmpdir, "test1")
	test2 := filepath.Join(tmpdir, "test2")

	ioutil.WriteFile(test1, []byte(""), 0644)
	// time.Sleep(0)
	ioutil.WriteFile(test2, []byte(""), 0644)

	info1, err := os.Stat(test1)
	require.NoError(t, err)
	info2, err := os.Stat(test2)
	require.NoError(t, err)

	require.NotEqual(t, info1.ModTime(), info2.ModTime())
}
