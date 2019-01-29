package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"
)

func TestProjectPathMissing(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir + "/nopenope")

	err := s.SetString("dummy", "")
	require.Error(t, err)
}

func TestInitialization(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	err := s.SetString("dummykey", "dummy")
	require.NoError(t, err)

	filet.DirContains(t, tmpdir, ".devbuddy/store")

	filet.DirContains(t, tmpdir, ".devbuddy/.gitignore")
	filet.FileSays(t, tmpdir+"/.devbuddy/.gitignore", []byte("*"))
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
		assert.Equal(t, testVal, val)
	}
}

func TestKeyEmpty(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	_, err := s.GetString("")
	require.EqualError(t, err, "empty string is not a valid key")

	err = s.SetString("", "")
	require.EqualError(t, err, "empty string is not a valid key")
}

func TestGetStringNotFound(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	s := New(tmpdir)

	val, err := s.GetString("doesnotexist")
	require.NoError(t, err)
	require.Equal(t, "", val)
}
