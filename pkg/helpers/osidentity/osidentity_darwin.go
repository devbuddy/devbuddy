package osidentity

// Detect returns an OS identifier.
func Detect() *Identity {
	return NewFromRuntime()
}
