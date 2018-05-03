package config

import (
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	homeDir      string
	DebugEnabled bool
	SourceDir    string
	dataDir      string
}

const defaultReleaseURL string = "https://api.github.com/repos/pior/dad/releases/latest"

func Load() (*Config, error) {
	homedir, err := getHomeDir()
	if err != nil {
		return nil, err
	}

	userDataDir := getXdgUserDataDir(homedir)

	c := Config{
		homeDir:      homedir,
		DebugEnabled: debugEnabled(),
		SourceDir:    filepath.Join(homedir, "src"),
		dataDir:      filepath.Join(userDataDir, "dad"),
	}
	return &c, nil
}

func debugEnabled() bool {
	return os.Getenv("DAD_DEBUG") != ""
}

func getHomeDir() (string, error) {
	home := os.Getenv("HOME")
	if home != "" {
		return home, nil
	}
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}

func getXdgUserDataDir(homedir string) string {
	dir := os.Getenv("XDG_DATA_HOME")
	if dir != "" {
		return dir
	}
	return filepath.Join(homedir, ".local/share")

}

func (c *Config) ReleaseURL() string {
	url := os.Getenv("DAD_RELEASE_URL")
	if url != "" {
		return url
	}
	return defaultReleaseURL
}

func (c *Config) HomeDir(elem ...string) string {
	elem = append([]string{c.homeDir}, elem...)
	return filepath.Join(elem...)
}

func (c *Config) DataDir(elem ...string) string {
	elem = append([]string{c.dataDir}, elem...)
	return filepath.Join(elem...)
}
