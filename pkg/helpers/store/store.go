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
	path string
}

// Open returns an instance of Store
func Open(projectPath, tableName string) (*Store, error) {
	path, err := projectmetadata.New(projectPath).Path()
	if err != nil {
		return nil, err
	}
	filename := "store-" + tableName
	return &Store{path: filepath.Join(path, filename)}, nil
}

func (s *Store) read() (*fileData, error) {
	data := &fileData{Entries: make(map[string]string)}

	if utils.PathExists(s.path) {
		serialized, err := ioutil.ReadFile(s.path)
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
	serialized, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.path, serialized, 0644)
}

// SetString stores a string for a key
func (s *Store) SetString(key, value string) error {
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
