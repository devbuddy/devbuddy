package utils

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/test"
	"github.com/stretchr/testify/require"
)

func Test_AppendOnlyFile(t *testing.T) {
	_, tmpfile := test.File(t, "testfile")

	err := AppendOnlyFile(tmpfile, []byte(""))
	require.Error(t, err)
	require.Contains(t, err.Error(), "no such file or directory")

	test.WriteFile(tmpfile, []byte("one"))

	err = AppendOnlyFile(tmpfile, []byte("two"))
	require.NoError(t, err)

	buffer := test.ReadFile(tmpfile)
	require.Equal(t, "onetwo", string(buffer))
}

func Test_WriteNewFile(t *testing.T) {
	_, tmpfile := test.File(t, "testfile")

	err := WriteNewFile(tmpfile, []byte("one"), 0664)
	require.NoError(t, err)

	content := test.ReadFile(tmpfile)
	require.Equal(t, "one", string(content))

	err = WriteNewFile(tmpfile, []byte("onetwo"), 0664)
	require.Error(t, err)
	require.Contains(t, err.Error(), "file exists")
}

func Test_WriteFile(t *testing.T) {
	_, tmpfile := test.File(t, "testfile")

	err := WriteFile(tmpfile, []byte("one"), 0664)
	require.NoError(t, err)

	content := test.ReadFile(tmpfile)
	require.Equal(t, "one", string(content))

	err = WriteFile(tmpfile, []byte("onetwo"), 0664)
	require.NoError(t, err)

	content = test.ReadFile(tmpfile)
	require.Equal(t, "onetwo", string(content))
}
