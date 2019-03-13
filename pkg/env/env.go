package env

import (
	"os"
	"strings"
)

// An Env provides a simple interface to manipulate environment variables
type Env struct {
	env         Variables
	verbatimEnv Variables
}

// New returns a new Env from an arbitrary list of variables
func New(env []string) *Env {
	return &Env{
		env:         NewVariables(env),
		verbatimEnv: NewVariables(env),
	}
}

// NewFromOS returns a new Env with variables from os.Environ()
func NewFromOS() (e *Env) {
	return New(os.Environ())
}

// Set adds or changes a variable
func (e *Env) Set(name string, value string) {
	e.env[name] = &Variable{name, value}
}

// Unset removes a variable if it exists
func (e *Env) Unset(name string) {
	delete(e.env, name)
}

// Get returns the value of a variable (defaults to empty string)
func (e *Env) Get(name string) string {
	variable := e.env[name]
	if variable != nil {
		return variable.Value
	}
	return ""
}

// Has returns whether the variable exists
func (e *Env) Has(name string) bool {
	return e.env[name] != nil
}

// Environ returns all variable as os.Environ() would
func (e *Env) Environ() []string {
	return e.env.AsEnviron()
}

// PrependToPath inserts a new path at the beginning of the PATH variable
func (e *Env) PrependToPath(path string) {
	elems := e.getPathParts()
	elems = append([]string{path}, elems...)
	e.setPathParts(elems...)
}

// RemoveFromPath removes all path entries matching a substring
func (e *Env) RemoveFromPath(substring string) {
	newElems := []string{}
	for _, elem := range e.getPathParts() {
		if !strings.Contains(elem, substring) {
			newElems = append(newElems, elem)
		}
	}
	e.setPathParts(newElems...)
}

// IsInPath returns true if one path of $PATH is exactly equal to the path provided
func (e *Env) IsInPath(path string) bool {
	for _, elem := range e.getPathParts() {
		if elem == path {
			return true
		}
	}
	return false
}

// Mutations returns a list of variable mutations (previous and current value)
func (e *Env) Mutations() (m []VariableMutation) {
	for _, current := range e.env {
		previous := e.verbatimEnv[current.Name]
		if !current.Eq(previous) {
			m = append(m, VariableMutation{Name: current.Name, Previous: previous, Current: current})
		}
	}
	for _, previous := range e.verbatimEnv {
		if _, present := e.env[previous.Name]; !present {
			m = append(m, VariableMutation{Name: previous.Name, Previous: previous, Current: nil})
		}
	}
	return m
}

func (e *Env) getPathParts() []string {
	return strings.Split(e.Get("PATH"), ":")
}

func (e *Env) setPathParts(elems ...string) {
	e.Set("PATH", strings.Join(elems, ":"))
}
