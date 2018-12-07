package osidentity

// Identity represent how your os behave
type Identity struct {
	platform string
	release  string
}

// IsDebianLike returns true if current platform behave like debian (including ubuntu)
func (i *Identity) IsDebianLike() bool {
	return i.platform == "linux" && i.release == "debian"
}

// IsMacOS returns true if current platform behave like macOS
func (i *Identity) IsMacOS() bool {
	return i.platform == "darwin"
}
