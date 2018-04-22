package env

import (
	"strings"
)

type Env struct {
	env         map[string]string
	verbatimEnv map[string]string
}

func New(env []string) (e *Env) {
	e = &Env{
		env:         make(map[string]string),
		verbatimEnv: make(map[string]string),
	}

	for _, pair := range env {
		parts := strings.SplitN(pair, "=", 2)
		e.env[parts[0]] = parts[1]
		e.verbatimEnv[parts[0]] = parts[1]
	}

	return
}

type VariableChange struct {
	Name    string
	Value   string
	Deleted bool
}

func (e *Env) Set(name string, value string) {
	e.env[name] = value
}

func (e *Env) Unset(name string) {
	delete(e.env, name)
}

func (e *Env) Get(name string) string {
	return e.env[name]
}

func (e *Env) getAndSplitPath() []string {
	return strings.Split(e.env["PATH"], ":")
}

func (e *Env) joinAndSetPath(elems ...string) {
	e.env["PATH"] = strings.Join(elems, ":")
}

func (e *Env) PrependToPath(path string) {
	elems := e.getAndSplitPath()
	elems = append([]string{path}, elems...)
	e.joinAndSetPath(elems...)
}

func (e *Env) RemoveFromPath(substring string) {
	newElems := []string{}
	for _, elem := range e.getAndSplitPath() {
		if !strings.Contains(elem, substring) {
			newElems = append(newElems, elem)
		}
	}
	e.joinAndSetPath(newElems...)
}

func (e *Env) Changed() []VariableChange {
	c := []VariableChange{}

	for k, v := range e.env {
		if v != e.verbatimEnv[k] {
			c = append(c, VariableChange{Name: k, Value: v})
		}
	}
	for k := range e.verbatimEnv {
		if _, ok := e.env[k]; !ok {
			c = append(c, VariableChange{Name: k, Deleted: true})
		}
	}
	return c
}
