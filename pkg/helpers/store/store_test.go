package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"
)

func open(path string) *Store {
	store, err := Open(path, "testTable")
	if err != nil {
		panic("failed to open the store")
	}
	return store
}

func TestSetGetString(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	testValues := []string{"DUMMY", "", "   "}

	for _, testVal := range testValues {
		err := open(tmpdir).SetString("key", testVal)
		require.NoError(t, err)

		val, err := open(tmpdir).GetString("key")
		require.NoError(t, err)
		assert.Equal(t, testVal, val)
	}
}

func TestKeyEmpty(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	_, err := open(tmpdir).GetString("")
	require.EqualError(t, err, "empty string is not a valid key")

	err = open(tmpdir).SetString("", "")
	require.EqualError(t, err, "empty string is not a valid key")
}

func TestGetStringNotFound(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	val, err := open(tmpdir).GetString("doesnotexist")
	require.NoError(t, err)
	require.Equal(t, "", val)
}
