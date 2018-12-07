package osidentity

import (
	"os"
	"runtime"
)

// Identity represent the os and the corresponding release.
type Identity struct {
	Platform string
	Release  string
}

// Detect returns an OS identifier.
func Detect() (*Identity, error) {
	variant := "unknown"

	if _, err := os.Stat("/etc/debian_version"); !os.IsNotExist(err) {
		variant = "debian"
	} else {
		return nil, err
	}

	return &Identity{runtime.GOOS, variant}, nil
}

// IsDebianLike returns true if current platform behave like debian (including ubuntu)
func (i *Identity) IsDebianLike() bool {
	return i.Platform == "linux" && i.Release == "debian"
}

// IsMacOS returns true if current platform behave like macOS
func (i *Identity) IsMacOS() bool {
	return i.Platform == "darwin"
}
