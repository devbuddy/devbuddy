package updatecheck

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTTL       = 24 * time.Hour
	installScriptURL = "https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh"
)

type Release struct {
	Version string
	URL     string
}

type Checker struct {
	CurrentVersion string
	CachePath      string
	Output         io.Writer
	Now            func() time.Time
	Latest         func() (*Release, error)
}

type Plan struct {
	Command string
	Note    string
}

type Pending struct {
	done    chan asyncResult
	printed bool
}

type cache struct {
	CheckedAt   time.Time `json:"checked_at"`
	Latest      Release   `json:"latest"`
	NotifiedFor string    `json:"notified_for"`
}

type asyncResult struct {
	currentVersion string
	latest         Release
	state          cache
	cachePath      string
	output         io.Writer
}

func (c Checker) Check() error {
	if !isReleaseVersion(c.CurrentVersion) {
		return nil
	}

	now := c.now()
	state := readCache(c.CachePath)
	latest := state.Latest
	shouldFetch := latest.Version == "" || now.Sub(state.CheckedAt) >= defaultTTL
	if shouldFetch {
		release, err := c.latest()
		if err != nil {
			return nil
		}
		latest = *release
		state.CheckedAt = now
		state.Latest = latest
	}

	if compareVersions(latest.Version, c.CurrentVersion) <= 0 {
		return writeCache(c.CachePath, state)
	}
	if state.NotifiedFor == latest.Version {
		return writeCache(c.CachePath, state)
	}

	writeNotice(c.output(), c.CurrentVersion, latest, UpgradePlan(c.CurrentVersion, latest.Version))
	state.NotifiedFor = latest.Version
	return writeCache(c.CachePath, state)
}

func (c Checker) Start() *Pending {
	if !isReleaseVersion(c.CurrentVersion) {
		return &Pending{}
	}

	now := c.now()
	state := readCache(c.CachePath)
	latest := state.Latest
	shouldFetch := latest.Version == "" || now.Sub(state.CheckedAt) >= defaultTTL
	if !shouldFetch {
		return &Pending{done: completedResult(c.noticeResult(latest, state))}
	}

	state.CheckedAt = now
	_ = writeCache(c.CachePath, state)

	pending := &Pending{done: make(chan asyncResult, 1)}
	go func() {
		release, err := c.latest()
		if err != nil {
			close(pending.done)
			return
		}
		state.Latest = *release
		if result, ok := c.noticeResult(*release, state); ok {
			pending.done <- result
		} else {
			_ = writeCache(c.CachePath, state)
		}
		close(pending.done)
	}()
	return pending
}

func (p *Pending) Finish() {
	if p == nil || p.done == nil || p.printed {
		return
	}
	select {
	case result, ok := <-p.done:
		if !ok {
			p.printed = true
			return
		}
		writeNotice(result.output, result.currentVersion, result.latest, UpgradePlan(result.currentVersion, result.latest.Version))
		result.state.NotifiedFor = result.latest.Version
		_ = writeCache(result.cachePath, result.state)
		p.printed = true
	default:
		return
	}
}

func (c Checker) noticeResult(latest Release, state cache) (asyncResult, bool) {
	if compareVersions(latest.Version, c.CurrentVersion) <= 0 {
		return asyncResult{}, false
	}
	if state.NotifiedFor == latest.Version {
		return asyncResult{}, false
	}
	return asyncResult{
		currentVersion: c.CurrentVersion,
		latest:         latest,
		state:          state,
		cachePath:      c.CachePath,
		output:         c.output(),
	}, true
}

func completedResult(result asyncResult, ok bool) chan asyncResult {
	done := make(chan asyncResult, 1)
	if ok {
		done <- result
	}
	close(done)
	return done
}

func ShouldRunForArgs(args []string) bool {
	return len(args) >= 2 && args[1] == "up"
}

func UpgradePlan(currentVersion string, latestVersion string) Plan {
	if ParseVersion(currentVersion).Source == "homebrew" {
		return Plan{Command: "brew upgrade devbuddy/devbuddy/devbuddy"}
	}
	return Plan{Command: fmt.Sprintf("curl -sSL %s | VERSION=%s sh", installScriptURL, latestVersion)}
}

func FetchLatestRelease(client *http.Client) (*Release, error) {
	if client == nil {
		client = http.DefaultClient
	}

	response, err := client.Get(releaseURL())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var payload struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.TagName == "" {
		return nil, fmt.Errorf("latest release response did not include tag_name")
	}
	return &Release{Version: payload.TagName, URL: payload.HTMLURL}, nil
}

func DefaultCachePath() string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "devbuddy", "update-check.json")
}

func New(currentVersion string, output io.Writer) Checker {
	client := &http.Client{Timeout: 2 * time.Second}
	return Checker{
		CurrentVersion: currentVersion,
		CachePath:      DefaultCachePath(),
		Output:         output,
		Now:            time.Now,
		Latest:         func() (*Release, error) { return FetchLatestRelease(client) },
	}
}

func (c Checker) now() time.Time {
	if c.Now != nil {
		return c.Now()
	}
	return time.Now()
}

func (c Checker) latest() (*Release, error) {
	if c.Latest != nil {
		return c.Latest()
	}
	return FetchLatestRelease(nil)
}

func (c Checker) output() io.Writer {
	if c.Output != nil {
		return c.Output
	}
	return io.Discard
}

func writeNotice(out io.Writer, current string, latest Release, plan Plan) {
	currentInfo := ParseVersion(current)
	fmt.Fprintf(out, "\nA new release of DevBuddy is available: %s -> %s\n", currentInfo.Version, latest.Version)
	if plan.Command != "" {
		fmt.Fprintf(out, "To upgrade, run: %s\n", plan.Command)
	} else if plan.Note != "" {
		fmt.Fprintf(out, "%s\n", plan.Note)
	}
	fmt.Fprintf(out, "%s\n", latest.URL)
}

func readCache(path string) cache {
	if path == "" {
		return cache{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cache{}
	}
	var state cache
	if err := json.Unmarshal(data, &state); err != nil {
		return cache{}
	}
	return state
}

func writeCache(path string, state cache) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func releaseURL() string {
	if url := os.Getenv("BUD_RELEASE_URL"); url != "" {
		return url
	}
	return "https://api.github.com/repos/devbuddy/devbuddy/releases/latest"
}

var versionRE = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)`)

type VersionInfo struct {
	Version string
	Source  string
}

func isReleaseVersion(version string) bool {
	return ParseVersion(version).Version != ""
}

func ParseVersion(version string) VersionInfo {
	match := versionRE.FindStringSubmatch(version)
	if match == nil {
		return VersionInfo{}
	}
	source := "manual"
	if strings.Contains(version, "-homebrew") {
		source = "homebrew"
	}
	return VersionInfo{Version: match[0], Source: source}
}

func compareVersions(left string, right string) int {
	lv, lok := parseVersion(ParseVersion(left).Version)
	rv, rok := parseVersion(ParseVersion(right).Version)
	if !lok || !rok {
		return 0
	}
	for i := range lv {
		if lv[i] > rv[i] {
			return 1
		}
		if lv[i] < rv[i] {
			return -1
		}
	}
	return 0
}

func parseVersion(version string) ([3]int, bool) {
	match := versionRE.FindStringSubmatch(version)
	if match == nil {
		return [3]int{}, false
	}
	var parsed [3]int
	for i := 0; i < 3; i++ {
		value, err := strconv.Atoi(match[i+1])
		if err != nil {
			return [3]int{}, false
		}
		parsed[i] = value
	}
	return parsed, true
}
