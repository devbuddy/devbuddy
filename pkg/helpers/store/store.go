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

// Key represents the key used to identify a stored
type Key string

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

func (s *Store) stateFilePath(kind string, key Key) string {
	return filepath.Join(s.path(), fmt.Sprintf("%s-%s", kind, key))
}

// KeyFromPath builds a Key for the path of a file in the project
func KeyFromPath(path string) Key {
	value := strings.Replace(path, string(filepath.Separator), "--", -1)
	return Key(value)
}

// Set stores an arbitrary value for a kind, key pair
func (s *Store) Set(kind string, key Key, value []byte) error {
	err := s.ensureInit()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.stateFilePath(kind, key), value, 0644)
}

// Set stores an arbitrary value for a kind, key pair
func (s *Store) SetString(kind string, key Key, value string) error {
	return s.Set(kind, key, []byte(value))
}

// Get retrieves an arbitrary value for a kind, key pair
func (s *Store) Get(kind string, key Key) ([]byte, error) {
	err := s.ensureInit()
	if err != nil {
		return nil, err
	}

	stateFilePath := s.stateFilePath(kind, key)

	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil
	}

	content, err := ioutil.ReadFile(stateFilePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// Get retrieves an arbitrary value for a kind, key pair
func (s *Store) GetString(kind string, key Key) (string, error) {
	value, err := s.Get(kind, key)
	return string(value), err
}

// RecordFileChange stores the modification time of a file.
func (s *Store) RecordFileChange(path string) error {
	checksum, err := utils.FileChecksum(filepath.Join(s.projectPath, path))
	if err != nil {
		return err
	}
	return s.Set("checksum", KeyFromPath(path), []byte(checksum))
}

// HasFileChanged detects whether a path has changed since the last call to RecordFileChange().
// Defaults to true if path doesn't exists or RecordFileChange() was never called.
func (s *Store) HasFileChanged(path string) (bool, error) {
	checksum, err := utils.FileChecksum(filepath.Join(s.projectPath, path))
	if err != nil {
		return true, nil
	}

	content, err := s.Get("checksum", KeyFromPath(path))
	if err != nil {
		return true, nil
	}

	return checksum != string(content), nil
}
