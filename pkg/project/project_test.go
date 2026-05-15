package project

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/stretchr/testify/require"
)

var cfg = &config.Config{SourceDir: "/src", DefaultPlatform: "github.com"}

func TestNewFromID(t *testing.T) {
	proj, err := NewFromID("golang/go", cfg)
	require.NoError(t, err, "NewFromID() failed")
	require.NotEqual(t, nil, proj)
	require.Equal(t, "go", proj.Name())
	require.Equal(t, "github.com:golang/go", proj.FullName())
	require.Equal(t, "/src/github.com/golang/go", proj.Path)

	url, err := proj.GetRemoteURL()
	require.NoError(t, err)
	require.Equal(t, "git@github.com:golang/go.git", url)

	require.Equal(t, "go-1999178051", proj.Slug())
}

func TestNewFromIDError(t *testing.T) {
	proj, err := NewFromID("", cfg)
	require.Error(t, err, "NewFromID() should fail")
	require.Nil(t, proj)
}

func TestNewFromID_DefaultOrg(t *testing.T) {
	cfg := &config.Config{SourceDir: "/src", DefaultPlatform: "github.com", DefaultOrg: "golang"}
	proj, err := NewFromID("go", cfg)
	require.NoError(t, err)
	require.Equal(t, "github.com:golang/go", proj.FullName())
}

func TestNewFromID_DefaultPlatform(t *testing.T) {
	cfg := &config.Config{SourceDir: "/src", DefaultPlatform: "gitlab.com"}
	proj, err := NewFromID("myorg/myrepo", cfg)
	require.NoError(t, err)
	require.Equal(t, "gitlab.com:myorg/myrepo", proj.FullName())
	require.Equal(t, "/src/gitlab.com/myorg/myrepo", proj.Path)

	url, err := proj.GetRemoteURL()
	require.NoError(t, err)
	require.Equal(t, "git@gitlab.com:myorg/myrepo.git", url)
}

func TestNewFromIDGithubFullURL(t *testing.T) {
	proj, err := NewFromID("git@github.com:golang/go.git", cfg)
	require.NoError(t, err, "NewFromID() failed")
	require.NotEqual(t, nil, proj)
	require.Equal(t, "go", proj.Name())
	require.Equal(t, "github.com:golang/go", proj.FullName())
	require.Equal(t, "/src/github.com/golang/go", proj.Path)

	url, err := proj.GetRemoteURL()
	require.Equal(t, "git@github.com:golang/go.git", url)

	require.Equal(t, "go-1999178051", proj.Slug())
}

func TestNewFromIDGithubFullHTTPURL(t *testing.T) {
	proj, err := NewFromID("https://github.com/golang/go.git", cfg)
	require.NoError(t, err, "NewFromID() failed")
	require.NotEqual(t, nil, proj)
	require.Equal(t, "go", proj.Name())
	require.Equal(t, "github.com:golang/go", proj.FullName())
	require.Equal(t, "/src/github.com/golang/go", proj.Path)

	url, err := proj.GetRemoteURL()
	require.Equal(t, "https://github.com/golang/go.git", url)

	require.Equal(t, "go-1999178051", proj.Slug())
}

func TestNewFromIDBitbucketFullURL(t *testing.T) {
	proj, err := NewFromID("git@bitbucket.org:zzzeek/dogpile.cache.git", cfg)
	require.NoError(t, err, "NewFromID() failed")
	require.NotEqual(t, nil, proj)
	require.Equal(t, "dogpile.cache", proj.Name())
	require.Equal(t, "bitbucket.org:zzzeek/dogpile.cache", proj.FullName())
	require.Equal(t, "/src/bitbucket.org/zzzeek/dogpile.cache", proj.Path)

	url, err := proj.GetRemoteURL()
	require.Equal(t, "git@bitbucket.org:zzzeek/dogpile.cache.git", url)

	require.Equal(t, "dogpile.cache-669781729", proj.Slug())
}

func TestNewFromIDCodeberg(t *testing.T) {
	proj, err := NewFromID("git@codeberg.org:myorg/myrepo.git", cfg)
	require.NoError(t, err)
	require.Equal(t, "codeberg.org:myorg/myrepo", proj.FullName())
	require.Equal(t, "/src/codeberg.org/myorg/myrepo", proj.Path)

	url, err := proj.GetRemoteURL()
	require.NoError(t, err)
	require.Equal(t, "git@codeberg.org:myorg/myrepo.git", url)
}

func TestNewFromIDSourcehut(t *testing.T) {
	proj, err := NewFromID("git@git.sr.ht:~myuser/myrepo", cfg)
	require.NoError(t, err)
	require.Equal(t, "git.sr.ht:~myuser/myrepo", proj.FullName())
	require.Equal(t, "/src/git.sr.ht/~myuser/myrepo", proj.Path)

	url, err := proj.GetRemoteURL()
	require.NoError(t, err)
	require.Equal(t, "git@git.sr.ht:~myuser/myrepo", url)
}

func TestNewFromIDHTTPSWithoutDotGit(t *testing.T) {
	proj, err := NewFromID("https://gitlab.com/myorg/myrepo", cfg)
	require.NoError(t, err)
	require.Equal(t, "gitlab.com:myorg/myrepo", proj.FullName())
	require.Equal(t, "/src/gitlab.com/myorg/myrepo", proj.Path)
}
