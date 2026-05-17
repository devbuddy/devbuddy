package release

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Kind string

const (
	KindMinor  Kind = "minor"
	KindPatch  Kind = "patch"
	KindRC     Kind = "rc"
	KindCustom Kind = "custom"
)

var versionRE = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)(?:-rc\.(\d+))?$`)

type version struct {
	major int
	minor int
	patch int
	rc    int
	isRC  bool
}

func NextVersion(tags []string, kind Kind, custom string) (string, error) {
	versions := parseTags(tags)
	if len(versions) == 0 {
		return "", fmt.Errorf("no release tags found")
	}

	slices.SortFunc(versions, compareVersions)
	latest := latestStable(versions)
	if latest == nil {
		return "", fmt.Errorf("no stable release tags found")
	}

	switch kind {
	case KindMinor:
		return latest.nextMinor().String(), nil
	case KindPatch:
		return latest.nextPatch().String(), nil
	case KindRC:
		return nextRC(versions, latest).String(), nil
	case KindCustom:
		return validateCustom(tags, custom)
	default:
		return "", fmt.Errorf("unknown release kind %q", kind)
	}
}

func parseTags(tags []string) []version {
	var versions []version
	for _, tag := range tags {
		v, ok := parseVersion(tag)
		if ok {
			versions = append(versions, v)
		}
	}
	return versions
}

func parseVersion(tag string) (version, bool) {
	matches := versionRE.FindStringSubmatch(strings.TrimSpace(tag))
	if matches == nil {
		return version{}, false
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	v := version{major: major, minor: minor, patch: patch}
	if matches[4] != "" {
		v.isRC = true
		v.rc, _ = strconv.Atoi(matches[4])
	}
	return v, true
}

func validateCustom(tags []string, custom string) (string, error) {
	if _, ok := parseVersion(custom); !ok {
		return "", fmt.Errorf("custom version must look like vMAJOR.MINOR.PATCH or vMAJOR.MINOR.PATCH-rc.N")
	}
	for _, tag := range tags {
		if strings.TrimSpace(tag) == custom {
			return "", fmt.Errorf("tag %s already exists", custom)
		}
	}
	return custom, nil
}

func latestStable(versions []version) *version {
	for i := len(versions) - 1; i >= 0; i-- {
		if !versions[i].isRC {
			return &versions[i]
		}
	}
	return nil
}

func nextRC(versions []version, latestStable *version) version {
	latest := versions[len(versions)-1]
	if latest.isRC {
		latest.rc++
		return latest
	}
	return latestStable.nextMinor().asRC()
}

func compareVersions(a, b version) int {
	if a.major != b.major {
		return a.major - b.major
	}
	if a.minor != b.minor {
		return a.minor - b.minor
	}
	if a.patch != b.patch {
		return a.patch - b.patch
	}
	if a.isRC != b.isRC {
		if a.isRC {
			return -1
		}
		return 1
	}
	return a.rc - b.rc
}

func (v version) nextMinor() version {
	return version{major: v.major, minor: v.minor + 1}
}

func (v version) nextPatch() version {
	return version{major: v.major, minor: v.minor, patch: v.patch + 1}
}

func (v version) asRC() version {
	v.isRC = true
	v.rc = 0
	return v
}

func (v version) String() string {
	result := fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
	if v.isRC {
		result += fmt.Sprintf("-rc.%d", v.rc)
	}
	return result
}
