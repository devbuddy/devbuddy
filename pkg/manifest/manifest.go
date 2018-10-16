package manifest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/devbuddy/devbuddy/pkg/utils"
)

var manifestFilename = "dev.yml"

// Manifest is a representation of the project manifest
type Manifest struct {
	Up       []interface{}       `yaml:"up"`
	Commands map[string]*Command `yaml:"commands"`
	Open     map[string]string   `yaml:"open"`
}

// Command is a representation of the `command` section of a manifest
type Command struct {
	Run         string `yaml:"run"`
	Description string `yaml:"desc"`
}

// Load returns a Manifest struct populated from a manifest file
func ExistsIn(path string) bool {
	return utils.PathExists(filepath.Join(path, manifestFilename))
}

// Load returns a Manifest struct populated from a manifest file
func Load(path string) (*Manifest, error) {
	manifestPath := filepath.Join(path, manifestFilename)
	if !utils.PathExists(manifestPath) {
		return nil, fmt.Errorf("no manifest at %s", manifestPath)
	}

	file, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	manifest := &Manifest{}
	err = yaml.Unmarshal(file, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}
