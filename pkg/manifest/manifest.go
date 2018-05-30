package manifest

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/pior/dad/pkg/utils"
)

var ManifestFilename = "dev.yml"

type Manifest struct {
	Path string
	Content
}

type Content struct {
	Up       []interface{}       `yaml:"up"`
	Commands map[string]*Command `yaml:"commands"`
	Open     map[string]string   `yaml:"open"`
}

type Command struct {
	Run         string `yaml:"run"`
	Description string `yaml:"desc"`
}

func Load(path string) (m *Manifest, err error) {
	manifestPath := filepath.Join(path, ManifestFilename)
	if !utils.PathExists(manifestPath) {
		return nil, nil
	}

	m = &Manifest{Path: manifestPath}

	file, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(file, &m.Content)
	if err != nil {
		return
	}

	return
}
