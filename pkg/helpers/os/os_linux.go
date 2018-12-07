package os

import (
	"os"
	"runtime"
)

// NewOS returns an OS identifier.
func NewOS() (*OS, error) {
	variant := "unknown"

	if _, err := os.Stat("/etc/debian_version"); !os.IsNotExist(err) {
		variant = "debian"
	} else {
		return nil, err
	}

	return &OS{runtime.GOOS, variant}, nil
}
