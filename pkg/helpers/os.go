package helpers

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
)

// OS represent the os and it's corresponding release.
type OS struct {
	version string
	release string
}

// NewOS returns an OS identifier.
func NewOS() (os *OS, err error) {
	if runtime.GOOS == "darwin" {
		osRelease, err := syscall.Sysctl("kern.osrelease")

		if err != nil {
			return nil, err
		}

		return &OS{"darwin", osRelease}, nil
	}

	return &OS{runtime.GOOS, ""}, nil
}

// GetVariant returns the variant of the os identified by `runtime`.
func (o *OS) GetVariant() (string, error) {
	if o.version == "darwin" {
		return o.getDarwinVariant()
	} else if o.version == "linux" {
		return o.getDarwinVariant()
	}

	return "", fmt.Errorf("Cannot identify %s", o.version)
}

func (o *OS) getDarwinVariant() (string, error) {
	versiomNumberList := strings.Split(o.release, ":")
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
		return "JAGUAR", nil
	} else if version == "5" {
		return "puma", nil
	}

	return "", fmt.Errorf("Cannot identify variant '%s' for darwin", version)
}

func (o *OS) getLinuxVariant() (string, error) {
	if _, err := os.Stat("/etc/debian_version"); !os.IsNotExist(err) {
		return "debian", nil
	}

	return "", fmt.Errorf("Cannot identify variant for linux")
}
