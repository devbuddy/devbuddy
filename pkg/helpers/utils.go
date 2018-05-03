package helpers

import (
	"io/ioutil"
	"os"
)

func makeTemporaryFile() (f *os.File, err error) {
	tmpFile, err := ioutil.TempFile("", dadPrefix)

	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}
