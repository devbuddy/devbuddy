package test

import (
	"os"
	"path"
	"testing"
)

// File creates a temporary directory and returns the full path of the file. Panic on error.
func File(tb testing.TB, filename string) (tmpdir, tmpfile string) {
	tb.Helper()

	tmpdir = tb.TempDir()
	tmpfile = path.Join(tmpdir, filename)
	return
}

// CreateFile creates a temporary directory, writes a file and returns the full path of the file. Panic on error.
func CreateFile(tb testing.TB, filename string, content []byte) (tmpdir, tmpfile string) {
	tb.Helper()

	tmpdir, tmpfile = File(tb, filename)

	err := os.WriteFile(tmpfile, content, 0600)
	if err != nil {
		panic("failed to write to " + tmpfile + ": " + err.Error())
	}

	return
}

// WriteFile writes, or overwrites a file. Panic on error.
func WriteFile(fullpath string, content []byte) {
	err := os.WriteFile(fullpath, content, 0600)
	if err != nil {
		panic("failed to write to " + fullpath + ": " + err.Error())
	}
}

// ReadFile returns the content of an existing file. Panic on error.
func ReadFile(fullpath string) []byte {
	content, err := os.ReadFile(fullpath)
	if err != nil {
		panic("failed to read from " + fullpath + ": " + err.Error())
	}
	return content
}
