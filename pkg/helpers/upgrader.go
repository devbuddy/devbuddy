package helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/termui"
)

type Upgrader struct {
	github  *Github
	client  *http.Client
	ui      *termui.UI
	useSudo bool
}

// NewUpgrader returns a new upgrade helper using the default http client
func NewUpgrader(cfg *config.Config, useSudo bool) (u *Upgrader) {
	return NewUpgraderWithHTTPClient(cfg, http.DefaultClient, useSudo)
}

// NewUpgraderWithHTTPClient returns a new upgrade helper using the provided http client
func NewUpgraderWithHTTPClient(cfg *config.Config, client *http.Client, useSudo bool) (u *Upgrader) {
	g := NewGithubWithClient(cfg, client)
	ui := termui.NewUI(cfg)

	return &Upgrader{
		github:  g,
		client:  client,
		useSudo: useSudo,
		ui:      ui,
	}
}

// Perform is fetching a new executable from `release`
//   and upgrading the executable at `destinationPath` with it
func (u *Upgrader) Perform(destinationPath string, sourceURL string) (err error) {
	data, err := u.github.Get(sourceURL)

	if err != nil {
		return
	}

	tmpFile, err := ioutil.TempFile("", "dad-")
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

	u.ui.CommandHeader(cmdline)

	return executor.NewShell(cmdline).Run()
}

// LatestRelease get latest release item for current platform
func (u *Upgrader) LatestRelease(plateform string) (release *GithubReleaseItem, err error) {
	release, err = u.github.LatestRelease(plateform)
	return
}
