package osidentity

import (
	"fmt"
	"os"
	"runtime"
)

// New returns an OS identifier.
func Detect() (*Identity, error) {
	variant := "unknown"

	if _, err := os.Stat("/etc/debian_version"); !os.IsNotExist(err) {
		variant = "debian"
	} else {
		return nil, err
	}

	return &Identity{runtime.GOOS, variant}, nil
}

// GetVariant returns the variant of the os identified by `runtime`.
func (i *Identity) GetVariant() (string, error) {
	if i.Release == "debian" {
		return "debian", nil
	}

	return "", fmt.Errorf("Cannot identify variant '%s' for linux", i.Release)
}
