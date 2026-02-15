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

func TestBuildGithubProjectURL(t *testing.T) {
	tests := []struct {
		name      string
		remoteURL string
		wantURL   string
	}{
		{
			name:      "SSH remote",
			remoteURL: "git@github.com:org1/repo1.git",
			wantURL:   "https://github.com/org1/repo1/tree/main",
		},
		{
			name:      "HTTPS remote with .git",
			remoteURL: "https://github.com/pior/runnable.git",
			wantURL:   "https://github.com/pior/runnable/tree/main",
		},
		{
			name:      "HTTPS remote without .git",
			remoteURL: "https://github.com/pior/runnable",
			wantURL:   "https://github.com/pior/runnable/tree/main",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			writer := test.Project(tmpdir)
			writer.CreateGitRepoWithRemote(t, tt.remoteURL)

			ctx := newTestContext()
			url, err := NewGitRepo(ctx, tmpdir).BuildGithubProjectURL()

			require.NoError(t, err)
			require.Equal(t, tt.wantURL, url)
		})
	}
}

func TestBuildGithubPullrequestURL(t *testing.T) {
	tests := []struct {
		name      string
		remoteURL string
		wantURL   string
	}{
		{
			name:      "SSH remote",
			remoteURL: "git@github.com:org1/repo1.git",
			wantURL:   "https://github.com/org1/repo1/pull/main?expand=1",
		},
		{
			name:      "HTTPS remote with .git",
			remoteURL: "https://github.com/pior/runnable.git",
			wantURL:   "https://github.com/pior/runnable/pull/main?expand=1",
		},
		{
			name:      "HTTPS remote without .git",
			remoteURL: "https://github.com/pior/runnable",
			wantURL:   "https://github.com/pior/runnable/pull/main?expand=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			writer := test.Project(tmpdir)
			writer.CreateGitRepoWithRemote(t, tt.remoteURL)

			ctx := newTestContext()
			url, err := NewGitRepo(ctx, tmpdir).BuildGithubPullrequestURL()

			require.NoError(t, err)
			require.Equal(t, tt.wantURL, url)
		})
	}
}

func TestBuildGithubURLUnrecognized(t *testing.T) {
	tmpdir := t.TempDir()
	writer := test.Project(tmpdir)
	writer.CreateGitRepoWithRemote(t, "https://gitlab.com/org/repo.git")

	ctx := newTestContext()
	_, err := NewGitRepo(ctx, tmpdir).BuildGithubProjectURL()

	require.Error(t, err)
	require.Contains(t, err.Error(), "unrecognized git remote url")
}
