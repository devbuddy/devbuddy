package termui

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/config"
)

type baseUI struct {
	out          *os.File
	debugEnabled bool
}

func newBaseUI() baseUI {
	return baseUI{
		out:          os.Stderr,
		debugEnabled: config.DebugEnabled(),
	}
}

func (u *baseUI) Debug(format string, params ...interface{}) {
	if u.debugEnabled {
		msg := fmt.Sprintf(format, params...)
		fmt.Fprintf(u.out, "DEBUG: %s\n", color.Gray(msg))
	}
}
