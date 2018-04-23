package utils

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"testing"

	filet "github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

func b64encode(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(msg))
}

func TestData(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	testfile := path.Join(tmpdir, "testfile")

	data := "Hello"

	err := DownloadFile(testfile, "https://httpbin.org/base64/"+b64encode(data))
	require.NoError(t, err)

	buffer, err := ioutil.ReadFile(testfile)
	require.NoError(t, err)

	require.Equal(t, data, string(buffer))
}

func TestStatus(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	testfile := path.Join(tmpdir, "testfile")

	err := DownloadFile(testfile, "https://httpbin.org/status/404")
	require.Error(t, err)
}

func TestRedirect(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	testfile := path.Join(tmpdir, "testfile")

	err := DownloadFile(testfile, "https://httpbin.org/redirect/2")
	require.NoError(t, err)

	fileInfo, err := os.Stat(testfile)
	require.NoError(t, err)

	require.NotZero(t, fileInfo.Size())
}
