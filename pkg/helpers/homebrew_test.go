package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsInCellar(t *testing.T) {
	prefix := t.TempDir()
	t.Setenv(EnvHomebrewPrefix, prefix)

	cellarPath := filepath.Join(prefix, "Cellar")
	caskroomPath := filepath.Join(prefix, "Caskroom")

	os.MkdirAll(filepath.Join(cellarPath, "curl", "1.2.3"), os.ModePerm)
	os.MkdirAll(filepath.Join(caskroomPath, "emacs", "26.1"), os.ModePerm)

	require.Truef(t, HomeBrewPackageIsInstalled("curl"), "Curl is properly installed in Cellar %s", cellarPath)
	require.Falsef(t, HomeBrewPackageIsInstalled("vim"), "vim is missing from Homebrew %s", prefix)
	require.True(t, HomeBrewPackageIsInstalled("emacs"), "Emacs is properly installed in caskroom %s", caskroomPath)
}
