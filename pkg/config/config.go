package config

import (
	"os"
	"path/filepath"
)

// Config represents the user system and environment configuration
type Config struct {
	homeDir      string // the user home directory
	DebugEnabled bool   // whether debug logging is activated
	SourceDir    string // projects base directory
	dataDir      string // directory managed by DevBuddy (languages distribs, virtualenvs...)
	DefaultOrg   string // [optional] default repo organisation
}

func NewTestConfig() *Config {
	return &Config{
		DebugEnabled: false,
	}
}

// Load returns a Config populated from the user environment
func Load() (*Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	userDataDir := getXdgUserDataDir(homedir)

	c := Config{
		homeDir:      homedir,
		DebugEnabled: debugEnabled(),
		SourceDir:    filepath.Join(homedir, "src"),
		dataDir:      filepath.Join(userDataDir, "bud"),
		DefaultOrg:   os.Getenv("BUD_DEFAULT_ORG"),
	}
	return &c, nil
}

func debugEnabled() bool {
	return os.Getenv("BUD_DEBUG") != ""
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
