package project

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/manifest"
)

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
		exists := manifest.ExistsIn(path)
		if exists {
			return NewFromPath(path), nil
		}

		// Continue searching in top directory
		path = filepath.Dir(path)
		if path == "/" {
			break
		}
	}

	return nil, ErrProjectNotFound
}
