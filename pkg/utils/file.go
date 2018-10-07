package utils

import (
	"bufio"
	"fmt"
	"hash/adler32"
	"io/ioutil"
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
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	checksum := fmt.Sprint(adler32.Checksum(content))
	return checksum, nil
}

func ReadLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return
}
