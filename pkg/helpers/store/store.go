package store

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/helpers/projectmetadata"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type fileData struct {
	Entries map[string]string `json:"entries"`
}

// Store is the place to record information or keep files about a project
type Store struct {
	projectMetadata *projectmetadata.ProjectMetadata
}

// New returns an instance of Store
func New(projectPath string) *Store {
	return &Store{projectMetadata: projectmetadata.New(projectPath)}
}

func (s *Store) path() (string, error) {
	path, err := s.projectMetadata.Path()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "store"), nil
}

func (s *Store) read() (*fileData, error) {
	path, err := s.path()
	if err != nil {
		return nil, err
	}

	data := &fileData{Entries: make(map[string]string)}

	if utils.PathExists(path) {
		serialized, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(serialized, data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *Store) write(data *fileData) error {
	path, err := s.path()
	if err != nil {
		return err
	}

	serialized, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, serialized, 0644)
}

// SetString stores a string for a key
func (s *Store) SetString(key string, value string) error {
	if key == "" {
		return errors.New("empty string is not a valid key")
	}

	data, err := s.read()
	if err != nil {
		return err
	}
	data.Entries[key] = value
	return s.write(data)
}

// GetString retrieves a string for a key
func (s *Store) GetString(key string) (string, error) {
	if key == "" {
		return "", errors.New("empty string is not a valid key")
	}

	data, err := s.read()
	if err != nil {
		return "", err
	}
	return data.Entries[key], nil
}
