package utils

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Flaque/filet"
)

func Test_AppendOnlyFile(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	nofile := path.Join(tmpdir, "testfile")

	err := AppendOnlyFile(nofile, []byte(""))
	require.Error(t, err)
	require.Contains(t, err.Error(), "no such file or directory")

	tmpfile := filet.TmpFile(t, tmpdir, "one").Name()

	err = AppendOnlyFile(tmpfile, []byte("two"))
	require.NoError(t, err)

	buffer, err := ioutil.ReadFile(tmpfile)
	require.NoError(t, err)
	require.Equal(t, "onetwo", string(buffer))
}

func Test_WriteNewFile(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	tmpfile := path.Join(tmpdir, "testfile")

	err := WriteNewFile(tmpfile, []byte("one"), 0664)
	require.NoError(t, err)

	buffer, err := ioutil.ReadFile(tmpfile)
	require.NoError(t, err)
	require.Equal(t, "one", string(buffer))

	err = WriteNewFile(tmpfile, []byte("onetwo"), 0664)
	require.Error(t, err)
	require.Contains(t, err.Error(), "file exists")
}

func Test_WriteFile(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")

	tmpfile := path.Join(tmpdir, "testfile")

	err := WriteFile(tmpfile, []byte("one"), 0664)
	require.NoError(t, err)

	buffer, err := ioutil.ReadFile(tmpfile)
	require.NoError(t, err)
	require.Equal(t, "one", string(buffer))

	err = WriteFile(tmpfile, []byte("onetwo"), 0664)
	require.NoError(t, err)

	buffer, err = ioutil.ReadFile(tmpfile)
	require.NoError(t, err)
	require.Equal(t, "onetwo", string(buffer))
}
