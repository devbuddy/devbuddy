package utils

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	filet "github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

type TestHandler struct{}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/download":
		w.Write([]byte("Hello"))
	case "/redirect":
		http.Redirect(w, r, "/download", 302)
	case "/loop":
		http.Redirect(w, r, "/loop", 302)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func TestData(t *testing.T) {
	server := httptest.NewServer(&TestHandler{})
	defer server.Close()

	tmpdir := filet.TmpDir(t, "")
	testfile := path.Join(tmpdir, "testfile")

	err := DownloadFile(testfile, server.URL+"/download")
	require.NoError(t, err)

	buffer, err := ioutil.ReadFile(testfile)
	require.NoError(t, err)

	require.Equal(t, "Hello", string(buffer))
}

func TestStatus(t *testing.T) {
	server := httptest.NewServer(&TestHandler{})
	defer server.Close()

	tmpdir := filet.TmpDir(t, "")
	testfile := path.Join(tmpdir, "testfile")

	err := DownloadFile(testfile, server.URL+"/notfound")
	require.Error(t, err)
}

func TestRedirect(t *testing.T) {
	server := httptest.NewServer(&TestHandler{})
	defer server.Close()

	tmpdir := filet.TmpDir(t, "")
	testfile := path.Join(tmpdir, "testfile")

	err := DownloadFile(testfile, server.URL+"/redirect")
	require.NoError(t, err)

	fileInfo, err := os.Stat(testfile)
	require.NoError(t, err)

	require.NotZero(t, fileInfo.Size())
}

func TestRedirectLimit(t *testing.T) {
	server := httptest.NewServer(&TestHandler{})
	defer server.Close()

	tmpdir := filet.TmpDir(t, "")
	testfile := path.Join(tmpdir, "testfile")

	err := DownloadFile(testfile, server.URL+"/loop")
	require.EqualErrorf(t, err, "Get /loop: stopped after 10 redirects", "")
}
