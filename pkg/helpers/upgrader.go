package helpers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
)

type Upgrader struct {
	ctx     *context.Context
	github  *Github
	client  *http.Client
	useSudo bool
}

// NewUpgrader returns a new upgrade helper using the default http client
func NewUpgrader(ctx *context.Context, useSudo bool) (u *Upgrader) {
	return NewUpgraderWithHTTPClient(ctx, http.DefaultClient, useSudo)
}

// NewUpgraderWithHTTPClient returns a new upgrade helper using the provided http client
func NewUpgraderWithHTTPClient(ctx *context.Context, client *http.Client, useSudo bool) (u *Upgrader) {
	g := NewGithubWithClient(client)

	return &Upgrader{
		ctx:     ctx,
		github:  g,
		client:  client,
		useSudo: useSudo,
	}
}

// Perform is fetching a new executable from `release`
// and upgrading the executable at `destinationPath` with it
func (u *Upgrader) Perform(destinationPath string, sourceURL string) (err error) {
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

	u.ctx.UI.CommandHeader(cmdline)

	return u.ctx.Executor.Run(executor.NewShell(cmdline)).Error
}

// LatestRelease get latest release item for current platform
func (u *Upgrader) LatestRelease(plateform string) (release *GithubReleaseItem, err error) {
	release, err = u.github.LatestRelease(plateform)
	return
}
