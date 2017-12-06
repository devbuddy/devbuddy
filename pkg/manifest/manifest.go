package manifest

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Path string
	Content
}

type Content struct {
	Up []map[interface{}]interface{} `yaml:"up"`
}

func Load(path string) (m *Manifest, err error) {
	m = &Manifest{Path: path}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(file, &m.Content)
	if err != nil {
		return
	}

	return
}
