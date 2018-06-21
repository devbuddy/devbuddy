package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsInCellar(t *testing.T) {
	prefix, err := ioutil.TempDir("/tmp", "dad-brew")
	require.NoError(t, err, "ioutil.TempDir() failed")

	cellarPath := filepath.Join(prefix, "Cellar")

	caskroomPath := filepath.Join(prefix, "Caskroom")

	os.MkdirAll(filepath.Join(cellarPath, "curl", "1.2.3"), os.ModePerm)
	os.MkdirAll(filepath.Join(caskroomPath, "emacs", "26.1"), os.ModePerm)

	h := NewHomebrewWithPrefix(prefix)

	require.Truef(t, h.IsInstalled("curl"), "Curl is properly installed in Cellar %s", cellarPath)
	require.Falsef(t, h.IsInstalled("vim"), "vim is missing from Homebrew %s", prefix)
	require.True(t, h.IsInstalled("emacs"), "Emacs is properly installed in caskroom %s", caskroomPath)
}
