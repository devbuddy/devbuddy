package config

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	SourceDir string
}

func Load() *Config {
	return &Config{
		SourceDir: ExpandDir("~/src"),
	}
}

func DebugEnabled() bool {
	return os.Getenv("DAD_DEBUG") != ""
}

func ExpandDir(path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(GetHomeDir(), path[2:])
	}
	return path
}

func GetHomeDir() string {
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}
	u, err := user.Current()
	if err != nil {
		panic("failed to determine the home dir")
	}
	return u.HomeDir
}

func PathExists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
