package utils

import (
	"os"
	"time"
)

func PathExists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Touch updates the atime and mtime of a file, after creating it if it doesn't exist.
func Touch(path string, atime, mtime time.Time) error {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	return os.Chtimes(path, atime, mtime)
}
