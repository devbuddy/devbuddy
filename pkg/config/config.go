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
	DataDir      string
}

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
		DataDir:      filepath.Join(userDataDir, "dad"),
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

func PathExists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
