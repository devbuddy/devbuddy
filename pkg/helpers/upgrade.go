package helpers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pior/dad/pkg/env"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
)

type Upgrade struct {
	github    *Github
	plateform string
	client    *http.Client
	useSudo   bool
}

// NewUpgrade returns a new upgrade helper using the default http client
func NewUpgrade(cfg *config.Config, useSudo bool) (u *Upgrade) {
	return NewUpgradeWithHTTPClient(cfg, http.DefaultClient, useSudo)
}

// NewUpgradeWithHTTPClient returns a new upgrade helper using the provided http client
func NewUpgradeWithHTTPClient(cfg *config.Config, client *http.Client, useSudo bool) (u *Upgrade) {
	g := NewGithubWithClient(cfg, client)
	env := env.NewFromOS()

	return &Upgrade{
		github:    g,
		plateform: env.Platform(),
		client:    client,
		useSudo:   useSudo,
	}
}

// Perform is fetching a new executable from `release` and upgrading the executable at `target` with it
func (u *Upgrade) Perform(target string, release *GithubReleaseItem) (err error) {
	data, err := release.Get(u.client)

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

	cmdline := u.buildCmdline(tmpFile.Name(), target)

	_, err = executor.NewShell(cmdline).Run()

	if err = tmpFile.Close(); err != nil {
		return
	}

	err = os.Remove(tmpFile.Name())

	return
}

func (u *Upgrade) buildCmdline(filename string, target string) string {
	if u.useSudo {
		return fmt.Sprintf("sudo cp %s %s", filename, target)
	}
	return fmt.Sprintf("cp %s %s", filename, target)
}

// LatestRelease get latest release item for current platform
func (u *Upgrade) LatestRelease() (release *GithubReleaseItem, err error) {
	release, err = u.LatestReleaseFor(u.plateform)

	return
}

// LatestReleaseFor get latest release item for `platform`
func (u *Upgrade) LatestReleaseFor(platform string) (release *GithubReleaseItem, err error) {
	release, err = u.github.LatestRelease(platform)

	return
}
