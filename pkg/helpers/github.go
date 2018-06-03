package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Github struct {
	client *http.Client
}

type GithubReleaseItem struct {
	Plateform   string `json:"name"`
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
		return nil, err
	}

	releases = &GithubReleaseList{}
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	if err = response.Body.Close(); err != nil {
		return
	}

	err = json.Unmarshal(body, releases)

	return
}

func releaseURL() string {
	url := os.Getenv("DAD_RELEASE_URL")
	if url != "" {
		return url
	}
	return "https://api.github.com/repos/pior/dad/releases/latest"
}

// LatestRelease get latest release url for a specific `platform`
func (g *Github) LatestRelease(plateform string) (*GithubReleaseItem, error) {
	releaseList, err := g.listReleases()

	if err != nil {
		return nil, err
	}

	for _, item := range releaseList.Items {
		if item.Plateform == plateform {
			return &item, nil
		}
	}

	err = fmt.Errorf("Cannot find release for %s", plateform)

	return nil, err
}

// Get download the content at `url`
func (g *Github) Get(url string) (data []byte, err error) {
	response, err := g.client.Get(url)

	if err != nil {
		return
	}

	data, err = ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	return data, response.Body.Close()
}
