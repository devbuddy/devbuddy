package termui

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"
)

type baseUI struct {
	out          *os.File
	debugEnabled bool
}

func (u *baseUI) Debug(format string, params ...interface{}) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		fmt.Fprintf(u.out, "BUD_DEBUG: %s\n", color.Gray(msg))
	}
}
