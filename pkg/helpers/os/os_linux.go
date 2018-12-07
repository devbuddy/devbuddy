package os

import (
	"fmt"
	"os"
	"runtime"
)

// New returns an OS identifier.
func New() (*OS, error) {
	variant := "unknown"

	if _, err := os.Stat("/etc/debian_version"); !os.IsNotExist(err) {
		variant = "debian"
	} else {
		return nil, err
	}

	return &OS{runtime.GOOS, variant}, nil
}

// GetVariant returns the variant of the os identified by `runtime`.
func (o *OS) GetVariant() (string, error) {
	if o.release == "debian" {
		return "debian", nil
	}

	return "", fmt.Errorf("Cannot identify variant '%s' for linux", o.release)
}
