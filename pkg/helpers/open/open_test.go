package open

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/test"
)

func newTestContext(proj *project.Project) *context.Context {
	return &context.Context{
		Cfg:      config.NewTestConfig(),
		Project:  proj,
		Env:      env.New([]string{}),
		Executor: executor.NewExecutor(),
	}
}

func TestFindLink(t *testing.T) {
	tmpdir := t.TempDir()
	writer := test.Project(tmpdir)
	writer.Manifest().WriteString("open: {doc: http://doc.com, logs: http://logs}")

	proj := project.NewFromPath(tmpdir)
	ctx := newTestContext(proj)

	_, err := FindLink(ctx, "")
	require.Error(t, err)

	_, err = FindLink(ctx, "unknown")
	require.Error(t, err)

	url, err := FindLink(ctx, "doc")
	require.NoError(t, err)
	require.Equal(t, "http://doc.com", url)
}

func TestFindLinkDefault(t *testing.T) {
	tmpdir := t.TempDir()
	writer := test.Project(tmpdir)
	writer.Manifest().WriteString("open: {doc: http://doc.com}")

	proj := project.NewFromPath(tmpdir)
	ctx := newTestContext(proj)

	url, err := FindLink(ctx, "")
	require.NoError(t, err)
	require.Equal(t, "http://doc.com", url)
}

func TestFindLinkGithub(t *testing.T) {
	tmpdir := t.TempDir()
	writer := test.Project(tmpdir)
	writer.CreateGitRepo(t)
	writer.Manifest().Empty()

	proj := project.NewFromPath(tmpdir)
	ctx := newTestContext(proj)

	nameToURL := map[string]string{
		"pullrequest": "https://github.com/org1/repo1/pull/main?expand=1",
		"pr":          "https://github.com/org1/repo1/pull/main?expand=1",
		"github":      "https://github.com/org1/repo1/tree/main",
		"gh":          "https://github.com/org1/repo1/tree/main",
	}
	for name, expectedURL := range nameToURL {
		url, err := FindLink(ctx, name)
		require.NoError(t, err)
		require.Equal(t, expectedURL, url)
	}
}

func TestPrintLinks(t *testing.T) {
	tmpdir := t.TempDir()
	writer := test.Project(tmpdir)

	proj := project.NewFromPath(tmpdir)

	writer.Manifest().WriteString("")

	err := PrintLinks(proj)
	require.Error(t, err)

	writer.Manifest().WriteString("open: {doc: http://doc.com}")
	err = PrintLinks(proj)
	require.NoError(t, err)
}
