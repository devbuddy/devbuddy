package utils

import (
	"fmt"
	"hash/adler32"
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

// FileChecksum reads a file and return the Adler32 checksum of its data as a string
func FileChecksum(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	checksum := fmt.Sprint(adler32.Checksum(content))
	return checksum, nil
}

// WriteNewFile writes data to a new file. Fails if the file exists.
func WriteNewFile(filename string, data []byte, perm os.FileMode) error {
	return write(filename, data, os.O_CREATE|os.O_WRONLY|os.O_EXCL, perm)
}

// WriteFile writes data to a new file, or rewrites an existing file.
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return write(filename, data, os.O_CREATE|os.O_WRONLY, perm)
}

// AppendOnlyFile appends to an existing file. Fails if the file does not exist.
func AppendOnlyFile(filename string, data []byte) error {
	return write(filename, data, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}

func write(filename string, data []byte, flag int, perm os.FileMode) error {
	file, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return err
	}

	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = file.Write(data)
	return err
}
