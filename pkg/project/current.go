package project

import (
	"errors"
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
			p := &Project{
				RepositoryName: filepath.Base(path),
				Path:           path,
				Manifest:       man,
			}
			return p, nil
		}

		path = filepath.Dir(path)
		if path == "/" {
			break
		}
	}

	return nil, ErrProjectNotFound
}
