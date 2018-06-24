package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

var dirName = ".devbuddy"

// Store is the place to record information or keep files about a project
type Store struct {
	projectPath string
}

// New returns an instance of Store
func New(projectPath string) *Store {
	return &Store{projectPath: projectPath}
}

func (s *Store) path() string {
	return filepath.Join(s.projectPath, dirName)
}

func (s *Store) ensureInit() (err error) {
	if !utils.PathExists(s.path()) {
		err = os.MkdirAll(s.path(), 0755)
		if err != nil {
			return
		}
	}

	gitignore := filepath.Join(s.path(), ".gitignore")
	if !utils.PathExists(gitignore) {
		err = ioutil.WriteFile(gitignore, []byte("*"), 0644)
		if err != nil {
			return
		}
	}
	return nil
}

func (s *Store) stateFilePath(kind, key string) string {
	return filepath.Join(s.path(), fmt.Sprintf("%s-%s", kind, key))
}

func makeKeyFromPath(path string) string {
	return strings.Replace(path, string(filepath.Separator), "--", -1)
}

// RecordFileChange stores the modification time of a file.
func (s *Store) RecordFileChange(path string) error {
	err := s.ensureInit()
	if err != nil {
		return err
	}

	info, err := os.Stat(filepath.Join(s.projectPath, path))
	if err != nil {
		return err
	}

	stateFilePath := s.stateFilePath("mtime", makeKeyFromPath(path))
	return utils.Touch(stateFilePath, info.ModTime(), info.ModTime())
}

// HasFileChanged detects whether a path has changed since the last call to RecordFileChange().
// Defaults to true if path doesn't exists or RecordFileChange() was never called.
func (s *Store) HasFileChanged(path string) bool {
	info, err := os.Stat(filepath.Join(s.projectPath, path))
	if err != nil {
		return true
	}

	stateInfo, err := os.Stat(s.stateFilePath("mtime", makeKeyFromPath(path)))
	if err != nil {
		return true
	}

	return info.ModTime().After(stateInfo.ModTime())
}
