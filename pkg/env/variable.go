package env

import (
	"fmt"
	"strings"
)

// Variable represents an environment variable
type Variable struct {
	Name  string
	Value string
}

func (v *Variable) Eq(other *Variable) bool {
	if v != nil && other != nil && v.Value != other.Value {
		return false
	}
	if (v == nil) != (other == nil) {
		return false
	}
	return true
}

type Variables map[string]*Variable

func NewVariables(environ []string) Variables {
	variables := map[string]*Variable{}

	for _, pair := range environ {
		parts := strings.SplitN(pair, "=", 2)
		variables[parts[0]] = &Variable{parts[0], parts[1]}
	}

	return variables
}

func (vs Variables) GetDefault(name string, defaultValue string) string {
	variable := vs[name]
	if variable != nil {
		return variable.Value
	}
	return defaultValue
}

func (vs Variables) Set(name, value string) {
	vs[name] = &Variable{name, value}
}

func (vs Variables) Unset(name string) {
	delete(vs, name)
}

func (vs Variables) AsEnviron() (vars []string) {
	for _, variable := range vs {
		vars = append(vars, variable.Name+"="+variable.Value)
	}
	return vars
}

// VariableMutation represents the change made on a variable
type VariableMutation struct {
	Name     string
	Previous *Variable
	Current  *Variable
}

func (m VariableMutation) DiffString() string {
	text := ""
	if m.Previous != nil {
		text += fmt.Sprintf("  - %s\n", m.Previous.Value)
	}
	if m.Current != nil {
		text += fmt.Sprintf("  + %s\n", m.Current.Value)
	}
	return text
}
