package os

import (
	"fmt"
	"strings"
)

// OS represent the os and it's corresponding release.
type OS struct {
	platform string
	release  string
}

func NewOSWithRelease(platform string, release string) *OS {
	return &OS{platform, release}
}

// GetVariant returns the variant of the os identified by `runtime`.
func (o *OS) GetVariant() (string, error) {
	if o.platform == "darwin" {
		return o.getDarwinVariant()
	} else if o.platform == "linux" {
		return o.getDarwinVariant()
	}

	return "", fmt.Errorf("Cannot identify %s", o.platform)
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
	if o.release == "debian" {
		return "debian", nil
	}

	return "", fmt.Errorf("Cannot identify variant for linux")
}
