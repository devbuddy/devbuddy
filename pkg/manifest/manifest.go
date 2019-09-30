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
	Env      map[string]string      `yaml:"env"`
	Up       []interface{}          `yaml:"up"`
	Commands map[string]interface{} `yaml:"commands"`
	Open     map[string]string      `yaml:"open"`
}

func (m *Manifest) commands() (map[string]*Command, error) {
	cmds := map[string]*Command{}

	for name, payload := range m.Commands {
		if p, ok := payload.(string); ok {
			cmds[name] = &Command{Run: p}
			continue
		}

		if p, ok := payload.(map[interface{}]interface{}); ok {
			if run, ok := p["run"]; ok {
				if runS, ok := run.(string); ok {
					cmds[name] = &Command{Run: runS}
					if desc, ok := p["desc"]; ok {
						if descS, ok := desc.(string); ok {
							cmds[name].Description = descS
						}
					}
					continue
				}
			}
		}

		return nil, fmt.Errorf("malformed manifest: invalid command %s", name)
	}

	return cmds, nil
}

func (m *Manifest) GetCommands() map[string]*Command {
	cmds, _ := m.commands()
	return cmds
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

	_, err = manifest.commands()
	if err != nil {
		return nil, err
	}

	return manifest, nil
}
