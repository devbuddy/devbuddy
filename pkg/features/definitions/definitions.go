package definitions

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type definition struct {
	Name       string
	Activate   func(string, *config.Config, *project.Project, *env.Env) error
	Refresh    func(string, *config.Config, *project.Project, *env.Env) error
	Deactivate func(string, *config.Config, *env.Env)
}

var definitions map[string]*definition

func init() {
	definitions = map[string]*definition{}
}

func Register(name string) *definition {
	if _, ok := definitions[name]; ok {
		panic(fmt.Sprint("Can't re-register a definition:", name))
	}
	definitions[name] = &definition{Name: name}
	return definitions[name]
}

func Get(name string) *definition {
	return definitions[name]
}

func Names() (names []string) {
	for name := range definitions {
		names = append(names, name)
	}
	return
}
