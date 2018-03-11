package project

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pior/dad/pkg/manifest"
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
		man, err := manifest.Load(path)
		if err != nil {
			return nil, err
		}
		if man != nil {
			p := &Project{
				RepositoryName: filepath.Base(path),
				Path:           path,
				Manifest:       man,
			}
			return p, nil
		}

		// Continue searching in top directory
		path = filepath.Dir(path)
		if path == "/" {
			break
		}
	}

	return nil, ErrProjectNotFound
}
