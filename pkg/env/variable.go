package env

import (
	"strings"
)

// Variable represents an environment variable
type variable struct {
	Name  string
	Value string
}

func (v *variable) eq(other *variable) bool {
	if v != nil && other != nil && v.Value != other.Value {
		return false
	}
	if (v == nil) != (other == nil) {
		return false
	}
	return true
}

type Variables map[string]*variable

func NewVariables(environ []string) Variables {
	variables := map[string]*variable{}

	for _, pair := range environ {
		parts := strings.SplitN(pair, "=", 2)
		variables[parts[0]] = &variable{parts[0], parts[1]}
	}

	return variables
}

func (vs Variables) getDefault(name string, defaultValue string) string {
	variable := vs[name]
	if variable != nil {
		return variable.Value
	}
	return defaultValue
}

func (vs Variables) set(name, value string) {
	vs[name] = &variable{name, value}
}

func (vs Variables) has(name string) bool {
	return vs[name] != nil
}

func (vs Variables) unset(name string) {
	delete(vs, name)
}

func (vs Variables) asEnviron() (vars []string) {
	for _, variable := range vs {
		vars = append(vars, variable.Name+"="+variable.Value)
	}
	return vars
}
