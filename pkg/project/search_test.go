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
	populateProjects(t, dir, "pior", []string{"dad", "george", "caravan", "pyramid_bugsnag"})
	populateProjects(t, dir, "george", []string{"carpe", "dorade", "gardon", "marlin"})

	tests := map[string]string{
		"marlin": "github.com:george/marlin",
		"mar":    "github.com:george/marlin",
		"rli":    "github.com:george/marlin",
		"mln":    "github.com:george/marlin",
		"gm":     "github.com:george/marlin",

		"car": "github.com:george/carpe", // multiple projects matches

		"george/carpe": "github.com:george/carpe", // with org name
		"george/car":   "github.com:george/carpe",
		"p/c":          "github.com:pior/caravan",
		"p/n":          "github.com:pior/caravan",

		"pyramid_bugsnag": "github.com:pior/pyramid_bugsnag", // with separator
		"pyramid":         "github.com:pior/pyramid_bugsnag",
		"_bug":            "github.com:pior/pyramid_bugsnag",
		"pb":              "github.com:pior/pyramid_bugsnag",
		"ppb":             "github.com:pior/pyramid_bugsnag",

		"george": "github.com:pior/george", // collision org<->project, project should win
		"gg":     "github.com:pior/george",
	}

	for input, expected := range tests {
		assert.Equal(t, expected, runFind(t, dir, input), "for input: %s", input)
	}
}

func TestNoMatch(t *testing.T) {
	defer filet.CleanUp(t)
	dir := filet.TmpDir(t, "")

	populateProjects(t, dir, "pior", []string{"dad"})

	cfg := &config.Config{SourceDir: dir}
	_, err := FindBestMatch("nope", cfg)
	require.Error(t, err, "FindBestMatch() should return an error when no project found")
	require.Equal(t, "no project found for nope", err.Error())
}
