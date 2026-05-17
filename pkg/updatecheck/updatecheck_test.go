package updatecheck

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCheckPrintsNoticeForNewReleaseAndCachesNotification(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "update.json")
	var out bytes.Buffer
	fetches := 0

	checker := Checker{
		CurrentVersion: "v0.16.1 [2026-05-17 12:00:00 +0000 UTC]",
		CachePath:      cachePath,
		Output:         &out,
		Now:            fixedTime,
		Latest: func() (*Release, error) {
			fetches++
			return &Release{Version: "v0.17.0", URL: "https://github.com/devbuddy/devbuddy/releases/tag/v0.17.0"}, nil
		},
	}

	err := checker.Check()

	require.NoError(t, err)
	require.Equal(t, 1, fetches)
	require.Contains(t, out.String(), "A new release of DevBuddy is available: v0.16.1 -> v0.17.0")
	require.Contains(t, out.String(), "VERSION=v0.17.0")
	require.Contains(t, out.String(), "https://github.com/devbuddy/devbuddy/releases/tag/v0.17.0")

	out.Reset()
	err = checker.Check()

	require.NoError(t, err)
	require.Equal(t, 1, fetches, "latest release should be cached until the TTL expires")
	require.Empty(t, out.String(), "notice should only be printed once for a latest version")
}

func TestAsyncCheckCachesImmediatelyAndPrintsNoticeOnlyAfterCompletion(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "update.json")
	var out bytes.Buffer
	releases := make(chan *Release, 1)

	checker := Checker{
		CurrentVersion: "v0.16.1",
		CachePath:      cachePath,
		Output:         &out,
		Now:            fixedTime,
		Latest: func() (*Release, error) {
			return <-releases, nil
		},
	}

	pending := checker.Start()

	cached := readCache(cachePath)
	require.Equal(t, fixedTime(), cached.CheckedAt)
	pending.Finish()
	require.Empty(t, out.String())

	releases <- &Release{Version: "v0.17.0", URL: "https://github.com/devbuddy/devbuddy/releases/tag/v0.17.0"}
	require.Eventually(t, func() bool {
		pending.Finish()
		return strings.Contains(out.String(), "A new release of DevBuddy is available")
	}, time.Second, 10*time.Millisecond)
}

func TestAsyncCheckDoesNotFetchAgainWhenCachedCheckIsFresh(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "update.json")
	require.NoError(t, writeCache(cachePath, cache{CheckedAt: fixedTime()}))
	fetches := 0

	checker := Checker{
		CurrentVersion: "v0.16.1",
		CachePath:      cachePath,
		Output:         &bytes.Buffer{},
		Now:            fixedTime,
		Latest: func() (*Release, error) {
			fetches++
			return &Release{Version: "v0.17.0", URL: "url"}, nil
		},
	}

	checker.Start().Finish()

	require.Equal(t, 0, fetches)
}

func TestCheckRefreshesAfterTTLButDoesNotRepeatSameNotice(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "update.json")
	var out bytes.Buffer
	fetches := 0
	now := fixedTime()

	checker := Checker{
		CurrentVersion: "v0.16.1",
		CachePath:      cachePath,
		Output:         &out,
		Now:            func() time.Time { return now },
		Latest: func() (*Release, error) {
			fetches++
			return &Release{Version: "v0.17.0", URL: "https://github.com/devbuddy/devbuddy/releases/tag/v0.17.0"}, nil
		},
	}
	require.NoError(t, checker.Check())
	out.Reset()

	now = now.Add(25 * time.Hour)
	require.NoError(t, checker.Check())

	require.Equal(t, 2, fetches)
	require.Empty(t, out.String())
}

func TestCheckUsesHomebrewUpgradeNoticeForHomebrewVersion(t *testing.T) {
	var out bytes.Buffer
	checker := Checker{
		CurrentVersion: "v0.16.1-homebrew",
		CachePath:      filepath.Join(t.TempDir(), "update.json"),
		Output:         &out,
		Now:            fixedTime,
		Latest:         func() (*Release, error) { return &Release{Version: "v0.17.0", URL: "url"}, nil },
	}

	require.NoError(t, checker.Check())

	require.Contains(t, out.String(), "To upgrade, run: brew upgrade devbuddy/devbuddy/devbuddy")
}

func TestCheckSilentlyIgnoresNetworkFailures(t *testing.T) {
	var out bytes.Buffer
	checker := Checker{
		CurrentVersion: "v0.16.1",
		CachePath:      filepath.Join(t.TempDir(), "update.json"),
		Output:         &out,
		Now:            fixedTime,
		Latest:         func() (*Release, error) { return nil, errors.New("offline") },
	}

	err := checker.Check()

	require.NoError(t, err)
	require.Empty(t, out.String())
}

func TestCheckSkipsUnparseableDevelopmentVersion(t *testing.T) {
	var out bytes.Buffer
	fetches := 0
	checker := Checker{
		CurrentVersion: "devel",
		CachePath:      filepath.Join(t.TempDir(), "update.json"),
		Output:         &out,
		Now:            fixedTime,
		Latest: func() (*Release, error) {
			fetches++
			return &Release{Version: "v0.17.0", URL: "https://github.com/devbuddy/devbuddy/releases/tag/v0.17.0"}, nil
		},
	}

	err := checker.Check()

	require.NoError(t, err)
	require.Equal(t, 0, fetches)
	require.Empty(t, out.String())
}

func TestCheckSkipsOlderOrSameRelease(t *testing.T) {
	for _, latest := range []string{"v0.16.1", "v0.15.9"} {
		t.Run(latest, func(t *testing.T) {
			var out bytes.Buffer
			checker := Checker{
				CurrentVersion: "v0.16.1",
				CachePath:      filepath.Join(t.TempDir(), "update.json"),
				Output:         &out,
				Now:            fixedTime,
				Latest:         func() (*Release, error) { return &Release{Version: latest, URL: "url"}, nil },
			}

			require.NoError(t, checker.Check())
			require.Empty(t, out.String())
		})
	}
}

func TestShouldRunForCommandSkipsShellHook(t *testing.T) {
	require.False(t, ShouldRunForArgs([]string{"bud", "--shell-hook"}))
	require.False(t, ShouldRunForArgs([]string{"bud", "--shell-init"}))
	require.False(t, ShouldRunForArgs([]string{"bud", "__complete", "up"}))
	require.False(t, ShouldRunForArgs([]string{"bud", "__completeNoDesc", "up"}))
	require.False(t, ShouldRunForArgs([]string{"bud", "upgrade"}))
	require.True(t, ShouldRunForArgs([]string{"bud", "up"}))
	require.True(t, ShouldRunForArgs([]string{"bud", "up", "--verbose"}))
	require.False(t, ShouldRunForArgs([]string{"bud", "--version"}))
	require.False(t, ShouldRunForArgs([]string{"bud", "inspect"}))
}

func TestUpgradePlanDetectsHomebrewFromVersionString(t *testing.T) {
	plan := UpgradePlan("v0.16.1-homebrew", "v0.17.0")

	require.Equal(t, "brew upgrade devbuddy/devbuddy/devbuddy", plan.Command)
	require.Equal(t, "", plan.Note)
}

func TestUpgradePlanUsesInstallScriptForReleaseBinary(t *testing.T) {
	plan := UpgradePlan("v0.16.1 [2026-05-17 12:00:00 +0000 UTC]", "v0.17.0")

	require.Equal(t, "curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | VERSION=v0.17.0 sh", plan.Command)
	require.Equal(t, "", plan.Note)
}

func TestUpgradePlanUsesInstallScriptForDevelopmentBuild(t *testing.T) {
	plan := UpgradePlan("dev-v0.16.1-3-gabcdef", "v0.17.0")

	require.Equal(t, "curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | VERSION=v0.17.0 sh", plan.Command)
	require.Equal(t, "", plan.Note)
}

func TestFetchLatestRelease(t *testing.T) {
	server := httpServer(t, `{"tag_name":"v0.17.0","html_url":"https://github.com/devbuddy/devbuddy/releases/tag/v0.17.0"}`)
	t.Setenv("BUD_RELEASE_URL", server.URL)

	release, err := FetchLatestRelease(http.DefaultClient)

	require.NoError(t, err)
	require.Equal(t, &Release{
		Version: "v0.17.0",
		URL:     "https://github.com/devbuddy/devbuddy/releases/tag/v0.17.0",
	}, release)
}

func fixedTime() time.Time {
	return time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
}

func httpServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
}
