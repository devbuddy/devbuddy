package os

// OS represent the os and it's corresponding release.
type OS struct {
	platform string
	release  string
}

func NewOSWithRelease(platform string, release string) *OS {
	return &OS{platform, release}
}
