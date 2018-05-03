package helpers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
)

type Upgrade struct {
	github    *Github
	plateform string
	client    *http.Client
	skipSudo  bool
}

func NewUpgrade(cfg *config.Config) (u *Upgrade) {
	return NewUpgradeWithHTTPClient(cfg, http.DefaultClient)
}

func NewUpgradeWithHTTPClient(cfg *config.Config, client *http.Client) (u *Upgrade) {
	g := NewGithubWithClient(cfg, client)
	plateform := cfg.Platform()

	return &Upgrade{
		github:    g,
		plateform: plateform,
		client:    client,
	}
}

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
	if u.skipSudo {
		return fmt.Sprintf("cp %s %s", filename, target)
	}

	return fmt.Sprintf("sudo cp %s %s", filename, target)
}

func (u *Upgrade) LatestRelease() (release *GithubReleaseItem, err error) {
	release, err = u.LatestReleaseFor(u.plateform)

	return
}

func (u *Upgrade) LatestReleaseFor(platform string) (release *GithubReleaseItem, err error) {
	release, err = u.github.LatestRelease(platform)

	return
}
