package manifest

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/pior/dad/pkg/utils"
)

var ManifestFilename = "dev.yml"

// Manifest is a representation of the project manifest
type Manifest struct {
	Up       []interface{}      `yaml:"up"`
	Commands map[string]Command `yaml:"commands"`
}

// Command is a representation of the `command` section of a manifest
type Command struct {
	Run         string `yaml:"run"`
	Description string `yaml:"desc"`
}

// Load returns a Manifest struct populated from a manifest file
func Load(path string) (m *Manifest, err error) {
	manifestPath := filepath.Join(path, ManifestFilename)
	if !utils.PathExists(manifestPath) {
		return nil, nil
	}

	file, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return
	}

	m = &Manifest{}
	err = yaml.Unmarshal(file, &m)
	if err != nil {
		return
	}

	return
}
