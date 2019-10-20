package helpers

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	filet "github.com/Flaque/filet"
)

func TestDebugLogWriter(t *testing.T) {
	defer filet.CleanUp(t)
	tmpdir := filet.TmpDir(t, "")
	testfile := path.Join(tmpdir, "debug.log")

	w := &DebugLogWriter{
		path:        testfile,
		maxFileSize: 90,
	}

	w.Write([]byte("aaaa\n"))
	w.Write([]byte("bbbb\n"))

	data, err := ioutil.ReadFile(testfile)
	require.NoError(t, err)
	require.Contains(t, string(data), "\naaaa\n")
	require.Contains(t, string(data), "\nbbbb\n")

	// log rotation
	w.Write([]byte("cccc\n"))

	data, err = ioutil.ReadFile(testfile + ".old")
	require.NoError(t, err)
	require.Contains(t, string(data), "\naaaa\n")
	require.Contains(t, string(data), "\nbbbb\n")

	data, err = ioutil.ReadFile(testfile)
	require.NoError(t, err)
	require.Contains(t, string(data), "\ncccc\n")

}
