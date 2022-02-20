package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type GOVersion struct {
	Major int
	Minor int
}

func NewGOVersion(major, minor int) GOVersion {
	return GOVersion{Major: major, Minor: minor}
}

func (v GOVersion) LessThan(o GOVersion) bool {
	if v.Major < o.Major {
		return true
	}

	if v.Major == o.Major {
		return v.Minor < o.Minor
	}

	return false
}

func ParseGOVersion(version string) (v GOVersion, err error) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)`)

	version = strings.TrimLeft(version, " ")

	matches := re.FindStringSubmatch(version)
	if len(matches) != 3 {
		return v, fmt.Errorf("not a valid go version")
	}

	v.Major, err = strconv.Atoi(matches[1])
	if err != nil {
		return v, fmt.Errorf("not a valid go version: %w", err)
	}
	v.Minor, err = strconv.Atoi(matches[2])
	if err != nil {
		return v, fmt.Errorf("not a valid go version: %w", err)
	}

	return v, nil
}
