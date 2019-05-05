package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/manifest"
)

// ErrProjectNotFound is returned when the project was not found
var ErrProjectNotFound = errors.New("project not found")

// FindCurrent returns the project that contains the current directory
func FindCurrent() (*Project, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error while searching for project: %s", err)
	}
	return findByPath(path)
}

func findByPath(path string) (*Project, error) {
	for {
		if manifest.ExistsIn(path) {
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
