package tasks

import (
	"fmt"

	color "github.com/logrusorgru/aurora"
)

type Unknown struct {
	definition interface{}
}

func (u *Unknown) Load(definition interface{}) (bool, error) {
	u.definition = definition
	return true, nil
}

func (u *Unknown) Perform() (err error) {
	fmt.Printf("%s %s: %+v\n", color.Brown("★"), color.Red("Unknown"), color.Brown(u.definition))
	return nil
}
