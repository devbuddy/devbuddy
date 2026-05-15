package project

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reShort = regexp.MustCompile(`^([\w.~-]+)/([\w.-]+)$`)
	reSSH   = regexp.MustCompile(`^git@([\w.-]+):([\w.~-]+)/([\w.-]+?)(?:\.git)?$`)
	reHTTPS = regexp.MustCompile(`^https://([\w.-]+)/([\w.~-]+)/([\w.-]+?)(?:\.git)?$`)
)

type hostingInfo struct {
	platform     string
	organisation string
	repository   string

	remoteURL string
}

func newHostingInfoByURL(url string, defaultPlatform string) (*hostingInfo, error) {
	url = strings.Trim(url, " ")

	if match := reShort.FindStringSubmatch(url); match != nil {
		return newHostingInfoFromParts("", defaultPlatform, match[1], match[2]), nil
	}

	if match := reSSH.FindStringSubmatch(url); match != nil {
		return newHostingInfoFromParts(url, match[1], match[2], match[3]), nil
	}

	if match := reHTTPS.FindStringSubmatch(url); match != nil {
		return newHostingInfoFromParts(url, match[1], match[2], match[3]), nil
	}

	return nil, fmt.Errorf("unrecognized project url: %s", url)
}

func newHostingInfoByPath(path string) *hostingInfo {
	return &hostingInfo{
		repository: filepath.Base(path),
	}
}

func newHostingInfoFromParts(remoteURL, platform, organisation, repository string) *hostingInfo {
	if remoteURL == "" {
		remoteURL = fmt.Sprintf("git@%s:%s/%s.git", platform, organisation, repository)
	}
	return &hostingInfo{
		platform:     platform,
		organisation: organisation,
		repository:   repository,
		remoteURL:    remoteURL,
	}
}
