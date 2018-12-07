package os

import (
	"runtime"
	"syscall"
)

// NewOS returns an OS identifier.
func NewOS() (*OS, error) {
	variant := "unknown"

	if variant, err := syscall.Sysctl("kern.osrelease"); err != nil {
		return nil, err
	}

	return &OS{runtime.GOOS, variant}, nil
}
