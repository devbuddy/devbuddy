package osidentity

import "runtime"

// Detect returns an OS identifier.
func Detect() *Identity {
	return &Identity{runtime.GOOS, ""}
}
