package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pior/dad/pkg/manifest"
)

var ManifestFilename = "dev.yml"

var ErrProjectNotFound = errors.New("project not found")

func FindCurrent() (*Project, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return findByPath(path)
}

func findByPath(path string) (*Project, error) {
	for {
		manifestPath := filepath.Join(path, ManifestFilename)
		if exists(manifestPath) {
			man, err := manifest.Load(manifestPath)
			if err != nil {
				return nil, err
			}
			return &Project{Path: path, Manifest: man}, nil
		}

		path = filepath.Dir(path)
		if path == "/" {
			break
		}
	}

	return nil, ErrProjectNotFound
}

// exists checks if a file or directory exists.
func exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return false
}
