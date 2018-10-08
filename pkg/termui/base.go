package termui

import (
	"fmt"
	"io"
	"log"
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
		Fprintf(u.out, "BUD_DEBUG: %s\n", color.Gray(msg))
	}
}

func Fprintf(w io.Writer, format string, a ...interface{}) {
	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		log.Fatalf("failed to write to console: %s", err)
	}
}
