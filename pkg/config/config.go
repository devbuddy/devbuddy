package config

import (
	"os"
	"os/user"
	"path/filepath"
)

// Config represents the user system and environment configuration
type Config struct {
	homeDir      string // the user home directory
	DebugEnabled bool   // whether debug logging is activated
	SourceDir    string // projects base directory
	dataDir      string // directory managed by Dad (languages distribs, virtualenvs...)
}

func NewTestConfig() *Config {
	return &Config{
		DebugEnabled: false,
	}
}

// Load returns a Config populated from the user environment
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

func (c *Config) HomeDir(elem ...string) string {
	elem = append([]string{c.homeDir}, elem...)
	return filepath.Join(elem...)
}

func (c *Config) DataDir(elem ...string) string {
	elem = append([]string{c.dataDir}, elem...)
	return filepath.Join(elem...)
}
