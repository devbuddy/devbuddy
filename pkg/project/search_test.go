package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func populateProjects(t *testing.T, sourceDir, org string, names []string) {
	for _, name := range names {
		err := os.MkdirAll(filepath.Join(sourceDir, "github.com", org, name), 0755)
		require.NoError(t, err, "failed to populate project")
	}
}

func runFind(t *testing.T, dir, input, defaultOrg string) string {
	cfg := &config.Config{SourceDir: dir, DefaultOrg: defaultOrg}
	proj, err := FindBestMatch(input, cfg)
	require.NoError(t, err, "FindBestMatch() crashed")
	return proj.FullName()
}

func TestMatches(t *testing.T) {
	dir := t.TempDir()
	populateProjects(t, dir, "pior", []string{"george", "caravan", "pyramid_bugsnag"})
	populateProjects(t, dir, "george", []string{"carpe", "dorade", "gardon", "marlin", "caravan"})

	tests := map[string]map[string]string{
		// No default organisation
		"": {
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
		},

		// With a default organization
		"george": {
			"carpe":   "github.com:george/carpe",
			"dorade":  "github.com:george/dorade",
			"gardon":  "github.com:george/gardon",
			"marlin":  "github.com:george/marlin",
			"caravan": "github.com:george/caravan",
		},
	}

	for defaultOrg, tt := range tests {
		for input, expected := range tt {
			assert.Equal(t, expected, runFind(t, dir, input, defaultOrg), "for input: %s", input)
		}
	}
}

func TestNoMatch(t *testing.T) {
	dir := t.TempDir()

	populateProjects(t, dir, "pior", []string{"whatever"})

	cfg := &config.Config{SourceDir: dir}
	_, err := FindBestMatch("nope", cfg)
	require.Error(t, err, "FindBestMatch() should return an error when no project found")
	require.Equal(t, "no project found for nope", err.Error())
}

func TestSearchingWithNoProject(t *testing.T) {
	dir := t.TempDir()

	cfg := &config.Config{SourceDir: dir}
	_, err := FindBestMatch("nope", cfg)
	require.Error(t, err, "FindBestMatch() should return an error when no project found")
	require.Equal(t, "no projects found at all! Try cloning one first", err.Error())
}

func TestFuzzySearch(t *testing.T) {
	index := []string{"ddg", "github", "heroku"}

	link := FindBestLinkMatch("github", index)
	require.Equal(t, "github", link)

	link = FindBestLinkMatch("dd", index)
	require.Equal(t, "ddg", link)

	link = FindBestLinkMatch("heru", index)
	require.Equal(t, "heroku", link)
}
