package utils

import (
	"os"
)

func PathExists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
