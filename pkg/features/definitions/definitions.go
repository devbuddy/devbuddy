package definitions

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type Definition struct {
	Name       string
	Activate   func(string, *config.Config, *project.Project, *env.Env) (bool, error)
	Deactivate func(string, *config.Config, *env.Env)
}

var definitions map[string]*Definition

func init() {
	definitions = map[string]*Definition{}
}

func Register(name string) *Definition {
	if _, ok := definitions[name]; ok {
		panic(fmt.Sprint("Can't re-register a definition:", name))
	}
	definitions[name] = &Definition{Name: name}
	return definitions[name]
}

func Get(name string) *Definition {
	return definitions[name]
}

func Names() (names []string) {
	for name := range definitions {
		names = append(names, name)
	}
	return
}
