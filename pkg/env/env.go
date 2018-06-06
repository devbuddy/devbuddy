package env

import (
	"os"
	"strings"
)

// An Env provides a simple interface to manipulate environment variables
type Env struct {
	env         map[string]string
	verbatimEnv map[string]string
}

// New returns a new Env from an arbitrary list of variables
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

// NewFromOS returns a new Env with variables from os.Environ()
func NewFromOS() (e *Env) {
	return New(os.Environ())
}

// A VariableChange represents the change made on a variable
type VariableChange struct {
	Name    string
	Value   string
	Deleted bool
}

// Set adds or changes a variable
func (e *Env) Set(name string, value string) {
	e.env[name] = value
}

// Unset removes a variable if it exists
func (e *Env) Unset(name string) {
	delete(e.env, name)
}

// Get returns the value of a variable (defaults to empty string)
func (e *Env) Get(name string) string {
	return e.env[name]
}

// Environ returns all variable as os.Environ() would
func (e *Env) Environ() (vars []string) {
	for name, value := range e.env {
		vars = append(vars, name+"="+value)
	}
	return vars
}

func (e *Env) getAndSplitPath() []string {
	return strings.Split(e.env["PATH"], ":")
}

func (e *Env) joinAndSetPath(elems ...string) {
	e.env["PATH"] = strings.Join(elems, ":")
}

// PrependToPath inserts a new path at the beginning of the PATH variable
func (e *Env) PrependToPath(path string) {
	elems := e.getAndSplitPath()
	elems = append([]string{path}, elems...)
	e.joinAndSetPath(elems...)
}

// RemoveFromPath removes all path entries matching a substring
func (e *Env) RemoveFromPath(substring string) {
	newElems := []string{}
	for _, elem := range e.getAndSplitPath() {
		if !strings.Contains(elem, substring) {
			newElems = append(newElems, elem)
		}
	}
	e.joinAndSetPath(newElems...)
}

// Changed returns all changes made on the variables
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
