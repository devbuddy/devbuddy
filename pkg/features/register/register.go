package register

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type Definition struct {
	Name       string
	Activate   func(string, *config.Config, *project.Project, *env.Env) error
	Deactivate func(string, *config.Config, *env.Env)
	Refresh    func(string, *config.Config, *project.Project, *env.Env) error
}

type Definitions struct {
	definitions map[string]*Definition
}

func New() *Definitions {
	return &Definitions{
		definitions: map[string]*Definition{},
	}
}

func (d *Definitions) Register(name string) *Definition {
	if _, ok := d.definitions[name]; ok {
		panic(fmt.Sprint("Can't re-register a definition:", name))
	}
	d.definitions[name] = &Definition{Name: name}
	return d.definitions[name]
}

func (d *Definitions) Get(name string) *Definition {
	return d.definitions[name]
}

func (d *Definitions) Names() (names []string) {
	for name := range d.definitions {
		names = append(names, name)
	}
	return
}
