package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func open(path string) *Store {
	store, err := Open(path, "testTable")
	if err != nil {
		panic("failed to open the store")
	}
	return store
}

func TestSetGetString(t *testing.T) {
	tmpdir := t.TempDir()

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
	tmpdir := t.TempDir()

	_, err := open(tmpdir).GetString("")
	require.EqualError(t, err, "empty string is not a valid key")

	err = open(tmpdir).SetString("", "")
	require.EqualError(t, err, "empty string is not a valid key")
}

func TestGetStringNotFound(t *testing.T) {
	tmpdir := t.TempDir()

	val, err := open(tmpdir).GetString("doesnotexist")
	require.NoError(t, err)
	require.Equal(t, "", val)
}
