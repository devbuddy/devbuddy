package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pior/dad/pkg/config"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func populateProjects(t *testing.T, sourceDir, org string, names []string) {
	for _, name := range names {
		err := os.MkdirAll(filepath.Join(sourceDir, "github.com", org, name), 0755)
		require.NoError(t, err, "failed to populate project")
	}
}

func runFind(t *testing.T, dir string, input string) string {
	cfg := &config.Config{SourceDir: dir}
	proj, err := FindBestMatch(input, cfg)
	require.NoError(t, err, "FindBestMatch() crashed")
	return proj.FullName()
}

func TestMatches(t *testing.T) {
	defer filet.CleanUp(t)
	dir := filet.TmpDir(t, "")
	populateProjects(t, dir, "george", []string{"carpe", "dorade", "gardon", "marlin"})
	populateProjects(t, dir, "pior", []string{"dad", "ecfg", "caravan", "pyramid_bugsnag"})

	tests := map[string]string{
		"marlin": "github.com:george/marlin",
		"mar":    "github.com:george/marlin",
		"rli":    "github.com:george/marlin",

		"pior/dad": "github.com:pior/dad",
		"pior/car": "github.com:pior/caravan",

		"pyramid_bugsnag": "github.com:pior/pyramid_bugsnag",
		"pyramid":         "github.com:pior/pyramid_bugsnag",
	}

	for input, expected := range tests {
		assert.Equal(t, expected, runFind(t, dir, input), "for input: %s", input)
	}
}
