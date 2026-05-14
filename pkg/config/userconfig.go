package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type UserConfig struct {
	Shell ShellConfig `yaml:"shell"`
}

type ShellConfig struct {
	DeferInit bool `yaml:"defer_init"`
}

func LoadUserConfig() *UserConfig {
	path := userConfigPath()
	if path == "" {
		return &UserConfig{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &UserConfig{}
	}

	var cfg UserConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return &UserConfig{}
	}
	return &cfg
}

func userConfigPath() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "devbuddy", "config.yml")
}
