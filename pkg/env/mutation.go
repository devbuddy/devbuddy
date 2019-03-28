package env

import "fmt"

// VariableMutation represents the change made on a variable
type VariableMutation struct {
	Name     string
	Previous *variable
	Current  *variable
}

// DiffString returns a representation of the mutation as a diff
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

func buildMutations(current, previous Variables) []VariableMutation {
	m := []VariableMutation{}

	for _, current := range current {
		previous := previous[current.Name]
		if !current.eq(previous) {
			m = append(m, VariableMutation{Name: current.Name, Previous: previous, Current: current})
		}
	}
	for _, previous := range previous {
		if _, present := current[previous.Name]; !present {
			m = append(m, VariableMutation{Name: previous.Name, Previous: previous, Current: nil})
		}
	}
	return m
}
