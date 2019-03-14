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
	variables := make(map[string]*Variable)

	for _, pair := range environ {
		parts := strings.SplitN(pair, "=", 2)
		variables[parts[0]] = &Variable{parts[0], parts[1]}
	}

	return variables
}

func (v Variables) AsEnviron() (vars []string) {
	for _, variable := range v {
		vars = append(vars, variable.Name+"="+variable.Value)
	}
	return vars
}

// A VariableChange represents the change made on a variable
type VariableChange struct {
	Name    string
	Value   string
	Deleted bool
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
