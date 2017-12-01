package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	SourceDir string
}

func Load() *Config {
	return &Config{
		SourceDir: filepath.Join(os.Getenv("HOME"), "src"),
	}
}
