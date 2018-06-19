package termui

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"
)

type baseUI struct {
	out          *os.File
	err          *os.File
	debugEnabled bool
}

func (u *baseUI) Debug(format string, params ...interface{}) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		fmt.Fprintf(u.err, "BUD_DEBUG: %s\n", color.Gray(msg))
	}
}
