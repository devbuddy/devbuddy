package colorfmt

import (
	"fmt"
	"io"
)

// Formatter interprets the color tags in a string and prints it to the file provided
type Formatter struct {
	file      io.Writer
	processor tagProcessor
}

// New returns a *Formatter
func New(file io.Writer, enableColor bool) *Formatter {
	f := &Formatter{file: file, processor: ignoreTagProcessor}
	if enableColor {
		f.processor = newColorizeTagProcessor().process
	}
	return f
}

// Printf interprets tags and print the colorized text
func (f *Formatter) Printf(format string, a ...interface{}) error {
	_, err := fmt.Fprintf(f.file, scan(format, '{', '}', f.processor), a...)
	return err
}
