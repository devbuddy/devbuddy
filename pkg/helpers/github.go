package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pior/dad/pkg/config"
)

const defaultReleaseURL string = "https://api.github.com/repos/pior/dad/releases/latest"

type Github struct {
	config *config.Config
	client *http.Client
}

type GithubReleaseItem struct {
	Plateform   string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	TagName     string `json:"tag_name"`
}
type GithubReleaseList struct {
	TagName string              `json:"tag_name"`
	Items   []GithubReleaseItem `json:"assets"`
}

func NewGithub(cfg *config.Config) (g *Github) {
	g = &Github{config: cfg, client: http.DefaultClient}

	return
}

func NewGithubWithClient(cfg *config.Config, client *http.Client) (g *Github) {
	g = &Github{config: cfg, client: client}

	return
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
	return defaultReleaseURL
}

// LatestRelease get latest release url for a specific `platform`
func (g *Github) LatestRelease(plateform string) (release *GithubReleaseItem, err error) {
	releaseList, err := g.listReleases()
	release = &GithubReleaseItem{}

	if err != nil {
		return
	}

	for _, *release = range releaseList.Items {
		if release.Plateform == plateform {
			release.TagName = releaseList.TagName
			return
		}
	}

	err = fmt.Errorf("Cannot find release for %s", plateform)
	release = nil

	return
}

func (item *GithubReleaseItem) Get(client *http.Client) (data []byte, err error) {
	response, err := client.Get(item.DownloadURL)

	if err != nil {
		return
	}

	data, err = ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	err = response.Body.Close()

	return
}
