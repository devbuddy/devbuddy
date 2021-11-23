package helpers

import (
	"log"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/termui"
	"github.com/devbuddy/devbuddy/pkg/test"
)

func TestUpgraderLatestRelease(t *testing.T) {
	_, tmpfile := test.File(t, "blob")

	r, err := recorder.New("fixtures/upgrader")
	require.NoError(t, err, "recorder.New() failed")

	defer func() {
		err = r.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()

	client := &http.Client{Transport: r}

	u := NewUpgraderWithHTTPClient(client, false)

	ui := termui.New(config.NewTestConfig())
	err = u.Perform(ui, tmpfile, "https://github.com/devbuddy/devbuddy/releases/download/v0.1.0/bud-darwin-amd64")
	require.NoError(t, err, "upgrader.Perform() failed")

	result := test.ReadFile(tmpfile)
	require.Equal(t, string(result), "Original data was too big")
}
