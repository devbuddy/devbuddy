package helpers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pior/dad/pkg/config"

	"github.com/stretchr/testify/require"

	"github.com/dnaeon/go-vcr/recorder"
)

const baseURL = "https://github.com/pior/dad/releases/download/v0.1.0"

func TestLatestRelease(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err, "config.Load() failed")

	r, err := recorder.New("fixtures/github")
	require.NoError(t, err, "recorder.New() failed")

	defer r.Stop()

	client := &http.Client{Transport: r}
	g := NewGithubWithClient(cfg, client)

	darwin, err := g.LatestRelease("dad-darwin-amd64")
	require.NoError(t, err, "github.LatestRelease() failed")

	require.Equal(t, fmt.Sprintf("%s/dad-darwin-amd64", baseURL), darwin.DownloadURL)

	linux, err := g.LatestRelease("dad-linux-amd64")
	require.NoError(t, err, "github.LatestRelease() failed")

	require.Equal(t, fmt.Sprintf("%s/dad-linux-amd64", baseURL), linux.DownloadURL)

	require.Equal(t, "dad-darwin-amd64", darwin.Plateform)
	require.Equal(t, "dad-linux-amd64", linux.Plateform)
}
