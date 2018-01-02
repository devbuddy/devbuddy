package tasks

import (
	"fmt"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/termui"
)

type Unknown struct {
	definition interface{}
}

func NewUnknown() Task {
	return &Unknown{}
}

func (u *Unknown) Load(definition interface{}) (bool, error) {
	u.definition = definition
	return true, nil
}

func (u *Unknown) Perform(ui *termui.UI) (err error) {
	fmt.Printf("%s %s: %+v\n", color.Brown("★"), color.Red("Unknown"), color.Brown(u.definition))
	return nil
}
