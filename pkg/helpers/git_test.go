package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/test"
)

func newTestContext() *context.Context {
	return &context.Context{
		Cfg:      config.NewTestConfig(),
		Env:      env.New([]string{}),
		Executor: executor.NewExecutor(),
	}
}

func TestGitGithubProjectURL(t *testing.T) {
	tmpdir := t.TempDir()
	writer := test.Project(tmpdir)
	writer.CreateGitRepo(t)

	ctx := newTestContext()
	url, err := NewGitRepo(ctx, tmpdir).BuildGithubProjectURL()

	require.NoError(t, err, "BuildGithubProjectURL() failed")
	require.Equal(t, "https://github.com/org1/repo1/tree/main", url)
}

func TestGitGithubPullrequestURL(t *testing.T) {
	tmpdir := t.TempDir()
	writer := test.Project(tmpdir)
	writer.CreateGitRepo(t)

	ctx := newTestContext()
	url, err := NewGitRepo(ctx, tmpdir).BuildGithubPullrequestURL()

	require.NoError(t, err, "BuildGithubProjectURL() failed")
	require.Equal(t, "https://github.com/org1/repo1/pull/main?expand=1", url)
}
