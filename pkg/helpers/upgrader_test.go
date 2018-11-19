package helpers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Flaque/filet"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

func TestUpgraderLatestRelease(t *testing.T) {
	defer filet.CleanUp(t)

	r, err := recorder.New("fixtures/upgrader")
	require.NoError(t, err, "recorder.New() failed")

	defer func() {
		err = r.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()

	target := filet.TmpFile(t, "", "")

	defer func() {
		err = os.Remove(target.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()

	client := &http.Client{Transport: r}

	u := NewUpgraderWithHTTPClient(client, false)

	ui := termui.New(config.NewTestConfig())
	err = u.Perform(ui, target.Name(), "https://github.com/devbuddy/devbuddy/releases/download/v0.1.0/bud-darwin-amd64")
	require.NoError(t, err, "upgrader.Perform() failed")

	result, err := ioutil.ReadFile(target.Name())
	require.NoError(t, err, "ioutil.ReadFile() failed")

	require.Equal(t, string(result), "Original data was too big")
}
