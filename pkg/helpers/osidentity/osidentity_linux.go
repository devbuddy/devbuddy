package osidentity

import (
	"runtime"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

// Detect returns an OS identifier.
func Detect() *Identity {
	variant := "unknown"

	if utils.PathExists("/etc/debian_version") {
		variant = "debian"
	}

	return &Identity{runtime.GOOS, variant}
}
