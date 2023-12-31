package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Github struct {
	client *http.Client
}

type GithubReleaseItem struct {
	Platform    string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}
type GithubReleaseList struct {
	TagName string              `json:"tag_name"`
	Items   []GithubReleaseItem `json:"assets"`
}

func NewGithub() *Github {
	return &Github{client: http.DefaultClient}
}

func NewGithubWithClient(client *http.Client) *Github {
	return &Github{client: client}
}

func (g *Github) listReleases() (releases *GithubReleaseList, err error) {
	response, err := g.client.Get(releaseURL())
	if err != nil {
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if err = response.Body.Close(); err != nil {
		return
	}

	releases = &GithubReleaseList{}
	err = json.Unmarshal(body, releases)
	return
}

func releaseURL() string {
	url := os.Getenv("BUD_RELEASE_URL")
	if url != "" {
		return url
	}
	return "https://api.github.com/repos/devbuddy/devbuddy/releases/latest"
}

// LatestRelease get latest release url for a specific `platform`
func (g *Github) LatestRelease(platform string) (*GithubReleaseItem, error) {
	releaseList, err := g.listReleases()
	if err != nil {
		return nil, err
	}

	for _, item := range releaseList.Items {
		if item.Platform == platform {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("Cannot find release for %s", platform)
}

// Get download the content at `url`
func (g *Github) Get(url string) (data []byte, err error) {
	response, err := g.client.Get(url)

	if err != nil {
		return
	}

	data, err = io.ReadAll(response.Body)

	if err != nil {
		return
	}

	return data, response.Body.Close()
}
