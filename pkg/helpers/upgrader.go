package helpers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
)

type Upgrader struct {
	github  *Github
	client  *http.Client
	useSudo bool
}

// NewUpgrader returns a new upgrade helper using the default http client
func NewUpgrader(cfg *config.Config, useSudo bool) (u *Upgrader) {
	return NewUpgraderWithHTTPClient(cfg, http.DefaultClient, useSudo)
}

// NewUpgraderWithHTTPClient returns a new upgrade helper using the provided http client
func NewUpgraderWithHTTPClient(cfg *config.Config, client *http.Client, useSudo bool) (u *Upgrader) {
	g := NewGithubWithClient(cfg, client)

	return &Upgrader{
		github:  g,
		client:  client,
		useSudo: useSudo,
	}
}

// Perform is fetching a new executable from `release`
//   and upgrading the executable at `destinationPath` with it
func (u *Upgrader) Perform(destinationPath string, sourceURL string) (err error) {
	data, err := u.github.Get(sourceURL)

	if err != nil {
		return
	}

	tmpFile, err := makeTemporaryFile()
	if err != nil {
		return
	}

	if _, err = tmpFile.Write(data); err != nil {
		return
	}

	cmdline := u.buildCmdline(tmpFile.Name(), destinationPath)

	if _, err = executor.NewShell(cmdline).Run(); err != nil {
		return
	}

	if err = tmpFile.Close(); err != nil {
		return
	}

	return os.Remove(tmpFile.Name())
}

func (u *Upgrader) buildCmdline(filename string, target string) string {
	if u.useSudo {
		return fmt.Sprintf("sudo cp %s %s", filename, target)
	}
	return fmt.Sprintf("cp %s %s", filename, target)
}

// LatestRelease get latest release item for current platform
func (u *Upgrader) LatestRelease(plateform string) (release *GithubReleaseItem, err error) {
	release, err = u.github.LatestRelease(plateform)

	return
}
