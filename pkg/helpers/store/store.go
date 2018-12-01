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

func (s *Store) pathForKey(key string) string {
	filePathForKey := strings.Replace(key, string(filepath.Separator), "--", -1)
	return filepath.Join(s.path(), filePathForKey)
}

func (s *Store) ensureInit() (err error) {
	if !utils.PathExists(s.projectPath) {
		return fmt.Errorf("failed to initialize the store: project path does not exist: %s", s.projectPath)
	}

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

// Set stores a byte slice for a key
func (s *Store) Set(key string, value []byte) error {
	err := s.ensureInit()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.pathForKey(key), value, 0644)
}

// SetString stores a string for a key
func (s *Store) SetString(key string, value string) error {
	return s.Set(key, []byte(value))
}

// Get retrieves a byte slice for a key
func (s *Store) Get(key string) ([]byte, error) {
	pathForKey := s.pathForKey(key)

	if _, err := os.Stat(pathForKey); os.IsNotExist(err) {
		return nil, nil
	}

	content, err := ioutil.ReadFile(pathForKey)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// GetString retrieves a string for a key
func (s *Store) GetString(key string) (string, error) {
	value, err := s.Get(key)
	return string(value), err
}

// DEPRECATED: don't use the RecordFileChange/HasFileChanged methods, they will be removed.
// It should not be part of the Store

// RecordFileChange stores the modification time of a file.
func (s *Store) RecordFileChange(path string) error {
	checksum, err := utils.FileChecksum(filepath.Join(s.projectPath, path))
	if err != nil {
		return err
	}
	return s.SetString(s.pathForKey("checksum-"+path), checksum)
}

// HasFileChanged detects whether a path has changed since the last call to RecordFileChange().
// Defaults to true if path doesn't exists or RecordFileChange() was never called.
func (s *Store) HasFileChanged(path string) (bool, error) {
	checksum, err := utils.FileChecksum(filepath.Join(s.projectPath, path))
	if err != nil {
		return true, nil
	}

	content, err := s.GetString(s.pathForKey("checksum-" + path))
	if err != nil {
		return true, nil
	}

	return checksum != content, nil
}
