package project

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/stretchr/testify/require"
)

var cfg = &config.Config{SourceDir: "/src"}

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

	proj, err = NewFromID("golang", cfg)
	require.Error(t, err, "NewFromID() should fail")
	require.Nil(t, proj)
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
