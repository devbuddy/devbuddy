package helpers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

type Upgrader struct {
	github  *Github
	client  *http.Client
	useSudo bool
}

// NewUpgrader returns a new upgrade helper using the default http client
func NewUpgrader(useSudo bool) (u *Upgrader) {
	return NewUpgraderWithHTTPClient(http.DefaultClient, useSudo)
}

// NewUpgraderWithHTTPClient returns a new upgrade helper using the provided http client
func NewUpgraderWithHTTPClient(client *http.Client, useSudo bool) (u *Upgrader) {
	g := NewGithubWithClient(client)

	return &Upgrader{
		github:  g,
		client:  client,
		useSudo: useSudo,
	}
}

// Perform is fetching a new executable from `release`
// and upgrading the executable at `destinationPath` with it
func (u *Upgrader) Perform(ui *termui.UI, destinationPath string, sourceURL string) (err error) {
	data, err := u.github.Get(sourceURL)

	if err != nil {
		return
	}

	tmpFile, err := os.CreateTemp("", "bud-")
	if err != nil {
		return
	}
	defer func() {
		err = tmpFile.Close()
		if err != nil {
			return
		}
		err = os.Remove(tmpFile.Name())
		if err != nil {
			return
		}
	}()

	if _, err = tmpFile.Write(data); err != nil {
		return
	}

	cmdline := fmt.Sprintf("cp %s %s", tmpFile.Name(), destinationPath)
	if u.useSudo {
		cmdline = fmt.Sprintf("sudo %s", cmdline)
	}

	ui.CommandHeader(cmdline)

	return executor.NewShell(cmdline).Run().Error
}

// LatestRelease get latest release item for current platform
func (u *Upgrader) LatestRelease(plateform string) (release *GithubReleaseItem, err error) {
	release, err = u.github.LatestRelease(plateform)
	return
}
