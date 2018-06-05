package helpers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnaeon/go-vcr/recorder"
)

const baseURL = "https://github.com/devbuddy/devbuddy/releases/download/v0.1.0"

func buildGithubClient(t *testing.T) (g *Github, r *recorder.Recorder) {
	r, err := recorder.New("fixtures/github")
	require.NoError(t, err, "recorder.New() failed")

	client := &http.Client{Transport: r}

	g = NewGithubWithClient(client)
	return
}
func TestLatestReleaseOnDarwin(t *testing.T) {
	g, r := buildGithubClient(t)

	darwin, err := g.LatestRelease("bud-darwin-amd64")
	require.NoError(t, err, "github.LatestRelease() failed")

	require.Equal(t, fmt.Sprintf("%s/bud-darwin-amd64", baseURL), darwin.DownloadURL)
	require.Equal(t, "bud-darwin-amd64", darwin.Plateform)

	err = r.Stop()
	require.NoError(t, err, "Recorder.Stop() failed")
}

func TestLatestReleaseOnLinux(t *testing.T) {
	g, r := buildGithubClient(t)

	linux, err := g.LatestRelease("bud-linux-amd64")
	require.NoError(t, err, "github.LatestRelease() failed")

	require.Equal(t, fmt.Sprintf("%s/bud-linux-amd64", baseURL), linux.DownloadURL)
	require.Equal(t, "bud-linux-amd64", linux.Plateform)

	err = r.Stop()
	require.NoError(t, err, "Recorder.Stop() failed")
}

func TestLatestReleaseOnUnknown(t *testing.T) {
	g, _ := buildGithubClient(t)

	_, err := g.LatestRelease("bud-linux-amd128")
	require.Error(t, err)
}
