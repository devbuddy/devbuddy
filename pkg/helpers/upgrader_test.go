package helpers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Flaque/filet"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/pior/dad/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestUpgraderLatestRelease(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err, "config.Load() failed")

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

	u := NewUpgraderWithHTTPClient(cfg, client, false)

	err = u.Perform(target.Name(), "https://github.com/pior/dad/releases/download/v0.1.0/dad-darwin-amd64")
	require.NoError(t, err, "upgrader.Perform() failed")

	result, err := ioutil.ReadFile(target.Name())
	require.NoError(t, err, "ioutil.ReadFile() failed")

	require.Equal(t, string(result), "Original data was too big")
}
