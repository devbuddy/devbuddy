package osidentity

import (
	"fmt"
	"os"
	"runtime"
	"github.com/cobaugh/osrelease"
)

// New returns an OS identifier.
func Detect() (*Identity, error) {
	release, err := osrelease.Read()

	if err != nil {
		return err
	}

	return DetectFromReleaseId(release['id']), nil
}

func DetectFromReleaseId(releaseId string) (*Identity, error) {
	return &Identity{runtime.GOOS, releaseId}, nil
}

// GetVariant returns the variant of the os identified by `runtime`.
func (i *Identity) GetVariant() (string, error) {
	if i.Release == "debian" {
		return "debian", nil
	}

	return "", fmt.Errorf("Cannot identify variant '%s' for linux", i.Release)
}
