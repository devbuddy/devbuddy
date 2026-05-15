package features

import (
	"errors"
	"os"
	"path"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	register.Register(ruby{})
}

type ruby struct{}

func (ruby) Name() string {
	return "ruby"
}

func (ruby) Activate(ctx *context.Context, param string) (bool, error) {
	warnRubyVersionMismatch(ctx, param)

	binPath := path.Join(helpers.RbEnvRoot(), "versions", param, "bin")
	if !utils.PathExists(binPath) {
		return true, nil
	}
	ctx.Env.PrependToPath(binPath)
	return false, nil
}

func (ruby) Deactivate(ctx *context.Context, param string) {}

// WatchedFiles re-activates the feature when .ruby-version changes so the
// mismatch warning stays in sync with the file's current contents.
func (ruby) WatchedFiles(param string) []string {
	return []string{".ruby-version"}
}

func warnRubyVersionMismatch(ctx *context.Context, activeVersion string) {
	if ctx.Project == nil {
		return
	}
	fileVersion, err := helpers.ReadRubyVersionFile(ctx.Project.Path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			ctx.UI.Warningf("ruby: failed to read .ruby-version: %s", err)
		}
		return
	}
	if fileVersion != activeVersion {
		ctx.UI.Warningf(
			"ruby: dev.yml requests %s but .ruby-version says %s. dev.yml wins; remove one to silence this warning.",
			activeVersion, fileVersion,
		)
	}
}
