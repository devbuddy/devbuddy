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
func NewFromOS() *Env) {
	return New(os.Environ())
}

// Set adds or changes a variable
func (e *Env) Set(name string, value string) {
	e.env.set(name, value)
}

// Unset removes a variable if it exists
func (e *Env) Unset(name string) {
	e.env.unset(name)
}

// Get returns the value of a variable (defaults to empty string)
func (e *Env) Get(name string) string {
	return e.env.getDefault(name, "")
}

// Has returns whether the variable exists
func (e *Env) Has(name string) bool {
	return e.env.has(name)
}

// Environ returns all variable as os.Environ() would
func (e *Env) Environ() []string {
	return e.env.asEnviron()
}

// PrependToPath inserts a new path at the beginning of the PATH variable
func (e *Env) PrependToPath(path string) {
	elems := e.getPathParts()
	elems = append([]string{path}, elems...)
	e.setPathParts(elems...)
}

func (e *Env) getPathParts() []string {
	return strings.Split(e.Get("PATH"), ":")
}

func (e *Env) setPathParts(elems ...string) {
	e.Set("PATH", strings.Join(elems, ":"))
}

// Mutations returns a list of variable mutations (previous and current value)
func (e *Env) Mutations() []VariableMutation {
	return buildMutations(e.env, e.verbatimEnv)
}
