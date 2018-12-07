package osidentity

import (
	"fmt"
	"runtime"
	"strings"
	"syscall"
)

// New returns an OS identifier.
func New() (i *Identity, err error) {
	variant := "unknown"

	if variant, err = syscall.Sysctl("kern.osrelease"); err != nil {
		return nil, err
	}

	return &Identity{runtime.GOOS, variant}, nil
}

// GetVariant returns the variant of the os identified by `runtime`.
func (i *Identity) GetVariant() (string, error) {
	versiomNumberList := strings.Split(i.Release, ":")
	versiomNumber := versiomNumberList[len(versiomNumberList)-1]
	versiomNumber = strings.TrimSpace(versiomNumber)

	version := strings.Split(versiomNumber, ".")[0]

	if version == "18" {
		return "mojave", nil
	} else if version == "17" {
		return "high sierra", nil
	} else if version == "16" {
		return "sierra", nil
	} else if version == "15" {
		return "el capitan", nil
	} else if version == "14" {
		return "yosemite", nil
	} else if version == "13" {
		return "mavericks", nil
	} else if version == "12" {
		return "mountain lion", nil
	} else if version == "11" {
		return "lion", nil
	} else if version == "10" {
		return "snow leopard", nil
	} else if version == "9" {
		return "leopard", nil
	} else if version == "8" {
		return "tiger", nil
	} else if version == "7" {
		return "panther", nil
	} else if version == "6" {
		return "jaguar", nil
	} else if version == "5" {
		return "puma", nil
	}

	return "", fmt.Errorf("Cannot identify variant '%s' for darwin", version)
}
