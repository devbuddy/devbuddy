package project

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func getPlatformNames() []string {
	return []string{"github.com", "bitbucket.org"}
}

type hostingInfo struct {
	platform     string
	organisation string
	repository   string

	remoteURL string
}

func newHostingInfoByURL(url string) (*hostingInfo, error) {
	url = strings.Trim(url, " ")

	reGithubFull := regexp.MustCompile(`^([\w.-]+)/([\w.-]+)$`)
	if match := reGithubFull.FindStringSubmatch(url); match != nil {
		return newGithubHostingInfo(match[1], match[2]), nil
	}

	reGithubGitURL := regexp.MustCompile(`^git@github.com:([\w.-]+)/([\w.-]+).git$`)
	if match := reGithubGitURL.FindStringSubmatch(url); match != nil {
		return newGithubHostingInfo(match[1], match[2]), nil
	}

	reBitbucketGitURL := regexp.MustCompile(`^git@bitbucket.org:([\w.-]+)/([\w.-]+).git$`)
	if match := reBitbucketGitURL.FindStringSubmatch(url); match != nil {
		return newBitbucketHostingInfo(match[1], match[2]), nil
	}

	return nil, fmt.Errorf("unrecognized project url: %s", url)
}

func newHostingInfoByPath(path string) *hostingInfo {
	return &hostingInfo{
		repository: filepath.Base(path),
	}
}

func newGithubHostingInfo(organisation, repository string) *hostingInfo {
	return &hostingInfo{
		platform:     "github.com",
		organisation: organisation,
		repository:   repository,
		remoteURL:    fmt.Sprintf("git@github.com:%s/%s.git", organisation, repository),
	}
}

func newBitbucketHostingInfo(organisation, repository string) *hostingInfo {
	return &hostingInfo{
		platform:     "bitbucket.org",
		organisation: organisation,
		repository:   repository,
		remoteURL:    fmt.Sprintf("git@bitbucket.org:%s/%s.git", organisation, repository),
	}
}
